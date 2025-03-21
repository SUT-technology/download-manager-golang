package downloadsrvc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/domain/model"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

// UpdateDownloadInDatabase updates (or appends) a download record and saves the JSON file.
func (d DownloadService) UpdateDownloadInDatabase(updated entity.Download) error {
	var (
		downloads []entity.Download
		err       error
	)
	ctx := context.Background()
	queryFunc := func(r *repository.Repo) error {
		downloads, err = r.Tables.Downloads.GetDownloads(ctx)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}

		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return err
	}
	found := false
	for i, d := range downloads {
		if d.ID == updated.ID {
			downloads[i] = updated
			found = true
			break
		}
	}
	if !found {
		downloads = append(downloads, updated)
	}

	queryFuncUpdate := func(r *repository.Repo) error {
		err = r.Tables.Downloads.UpdateDownloads(ctx, downloads)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}

		return nil
	}

	if err = d.db.Query(queryFuncUpdate); err != nil {
		fmt.Println("Error saving database:", err)
	}

	return nil
}

func (d DownloadService) DownloadWorker(download *entity.Download, progressChan chan<- model.DownloadProgressMsg, controlChan <-chan model.DownloadControlMessage) {
	// Mark the download as starting.
	download.Status = "downloading"
	d.UpdateDownloadInDatabase(*download)
	incrementActiveDownload(download.QueueId, download.ID)

	const chunkSize = 32 * 1024 // 32 KB per iteration.
	lastBytes := download.Downloaded
    lastUpdate := time.Now()
	for download.Downloaded < download.TotalSize {
		// Check for pause/resume commands.
		// Look up the queue to retrieve MaximumBandWidth.
		q, err := d.GetQueueById(download.QueueId)
		if err != nil {
			fmt.Printf("Queue not found for download %d\n", download.ID)
			break
		}

		waitForAllowedTime(*q)
		select {
		case ctrl := <-controlChan:
			if ctrl == model.PauseCommand {
				download.Status = "paused"
				progressChan <- model.DownloadProgressMsg{
					DownloadID: download.ID,
					Progress:   float64(download.Downloaded) / float64(download.TotalSize) * 100,
					Speed:      0,
					Status:     download.Status,
					Downloaded: download.Downloaded,
				}
				d.UpdateDownloadInDatabase(*download)
				// Decrement active count when pausedownload.
				decrementActiveDownload(download.QueueId, download.ID)
				// Block until resume is receivedownload.
				for {
                    ctrl2 := <-controlChan
                    if ctrl2 == model.ResumeCommand {
                        download.Status = "downloading"
                        lastUpdate = time.Now()
                        lastBytes = download.Downloaded
                        break
                    }
                }
			}
		default:
			// Set up a GET request with the appropriate Range header.
            req, err := http.NewRequest("GET", download.URL, nil)
            if err != nil {
                fmt.Println("Error creating GET request:", err)
                break
            }
            rangeHeader := fmt.Sprintf("bytes=%d-%d", download.Downloaded, download.Downloaded+chunkSize-1)
            if download.Downloaded+chunkSize > download.TotalSize {
                rangeHeader = fmt.Sprintf("bytes=%d-", download.Downloaded)
            }
            req.Header.Set("Range", rangeHeader)

            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                fmt.Println("Error performing GET request:", err)
                break
            }

            // Open file in append mode.
            outFile, err := os.OpenFile(download.OutPath, os.O_WRONLY|os.O_APPEND, 0644)
            if err != nil {
                fmt.Println("Error opening output file:", err)
                resp.Body.Close()
                break
            }

            // Copy up to chunkSize bytes from the response.
            n, err := io.CopyN(outFile, resp.Body, chunkSize)
            outFile.Close()
            resp.Body.Close()

            if n > 0 {
                download.Downloaded += n
                if download.Downloaded > download.TotalSize {
                    download.Downloaded = download.TotalSize
                }
                download.Progress = float64(download.Downloaded) / float64(download.TotalSize) * 100

                now := time.Now()
                elapsed := now.Sub(lastUpdate).Seconds()
                var speed float64
                if elapsed > 0 {
                    speed = float64(download.Downloaded-lastBytes) / elapsed / 1024.0
                }
                download.CurrentSpeed = speed

                progressChan <- model.DownloadProgressMsg{
                    DownloadID: download.ID,
                    Progress:   download.Progress,
                    Speed:      download.CurrentSpeed,
                    Status:     download.Status,
                    Downloaded: download.Downloaded,
                }
                d.UpdateDownloadInDatabase(*download)

                lastUpdate = now
                lastBytes = download.Downloaded
            }
            if err != nil {
                // Ignore io.EOF as it signals the end of the current request.
                if err != io.EOF {
                    fmt.Println("Error reading response body:", err)
                }
            }
		}

	}

	download.Status = "completed"
	progressChan <- model.DownloadProgressMsg{
		DownloadID: download.ID,
		Progress:   100,
		Speed:      0,
		Status:     download.Status,
		Downloaded: download.Downloaded,
	}
	d.UpdateDownloadInDatabase(*download)
	decrementActiveDownload(download.QueueId, download.ID)
}

func (d DownloadService) StartPendingDownloads(downloads []entity.Download) {
	for i := range downloads {
		down := &downloads[i]
		if down.Status == "pending" || down.Status == "paused" || down.Status == "downloading" {
			controlChan := make(chan model.DownloadControlMessage)
			ControlChannels[down.ID] = controlChan
			go d.DownloadWorker(down, ProgressChan, controlChan)
		}
	}
}
