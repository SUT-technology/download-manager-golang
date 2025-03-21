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

    for download.Downloaded < download.TotalSize {
        // Retrieve the queue which holds the bandwidth limit.
        q, err := d.GetQueueById(download.QueueId)
        if err != nil {
            fmt.Printf("Queue not found for download %d\n", download.ID)
            break
        }

        // Wait until the current time is allowed by the queue’s time window.
        waitForAllowedTime(*q)

        // Check for pause/resume commands.
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
                decrementActiveDownload(download.QueueId, download.ID)
                // Wait until resume is received.
                for {
                    ctrl2 := <-controlChan
                    if ctrl2 == model.ResumeCommand {
                        download.Status = "downloading"
                        d.UpdateDownloadInDatabase(*download)
                        incrementActiveDownload(download.QueueId, download.ID)
                        break
                    }
                }
            }
        default:
            // No control command, proceed.
        }

        startChunk := time.Now()

        // Prepare the request.
        req, err := http.NewRequest("GET", download.URL, nil)
        if err != nil {
            fmt.Println("Error creating GET request:", err)
            break
        }
        endRange := download.Downloaded + chunkSize - 1
        if download.Downloaded+chunkSize > download.TotalSize {
            endRange = download.TotalSize - 1
        }
        req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", download.Downloaded, endRange))

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error performing GET request:", err)
            break
        }

        // If the server does not support Range requests:
        if resp.StatusCode == http.StatusOK {
            // If we have not yet downloaded anything, then assume a full download.
            if download.Downloaded == 0 {
                fmt.Println("Server did not honor Range, performing full download")
                // Overwrite the file; use os.Create (which truncates).
                outFile, err := os.Create(download.OutPath)
                if err != nil {
                    fmt.Println("Error creating output file for full download:", err)
                    resp.Body.Close()
                    break
                }
                n, err := io.Copy(outFile, resp.Body)
                outFile.Close()
                resp.Body.Close()
                download.Downloaded = n
                download.Progress = 100
                download.CurrentSpeed = float64(n) / (time.Since(startChunk).Seconds() * 1024) // KB/s
                progressChan <- model.DownloadProgressMsg{
                    DownloadID: download.ID,
                    Progress:   download.Progress,
                    Speed:      download.CurrentSpeed,
                    Status:     download.Status,
                    Downloaded: download.Downloaded,
                }
                d.UpdateDownloadInDatabase(*download)
                break
            } else {
                // If we already have some bytes, then the server does not support resume.
                fmt.Println("Server does not support resume")
                resp.Body.Close()
                break
            }
        } else if resp.StatusCode != http.StatusPartialContent {
            // Unexpected status code.
            fmt.Printf("Unexpected response status %s for download %d\n", resp.Status, download.ID)
            resp.Body.Close()
            break
        }

        // Expected status is 206; continue to read a chunk.
        outFile, err := os.OpenFile(download.OutPath, os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            fmt.Println("Error opening output file:", err)
            resp.Body.Close()
            break
        }

        n, err := io.CopyN(outFile, resp.Body, chunkSize)
        outFile.Close()
        resp.Body.Close()

        // Calculate time taken to download the chunk.
        elapsed := time.Since(startChunk).Seconds()

        activeCount := GetActiveDownload(download.QueueId)
        if activeCount < 1 {
            activeCount = 1
        }
        allocatedSpeed := q.MaximumBandWidth / float64(activeCount) // in KB/s

        chunkKB := float64(n) / 1024.0
        desiredDuration := chunkKB / allocatedSpeed
        sleepTime := desiredDuration - elapsed
        if sleepTime > 0 {
            time.Sleep(time.Duration(sleepTime * float64(time.Second)))
        }

        effectiveElapsed := elapsed
        if sleepTime > 0 {
            effectiveElapsed += sleepTime
        }
        measuredSpeed := chunkKB / effectiveElapsed

        download.Downloaded += n
        if download.Downloaded > download.TotalSize {
            download.Downloaded = download.TotalSize
        }
        download.Progress = float64(download.Downloaded) / float64(download.TotalSize) * 100
        download.CurrentSpeed = measuredSpeed

        progressChan <- model.DownloadProgressMsg{
            DownloadID: download.ID,
            Progress:   download.Progress,
            Speed:      measuredSpeed,
            Status:     download.Status,
            Downloaded: download.Downloaded,
        }
        d.UpdateDownloadInDatabase(*download)

        if err != nil && err != io.EOF {
            fmt.Println("Error reading chunk:", err)
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
