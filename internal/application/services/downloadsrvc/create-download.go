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

	download.Status = "downloading"
	d.UpdateDownloadInDatabase(*download)
	incrementActiveDownload(download.QueueId, download.ID)

	const chunkSize = 32 * 1024

	outFile, err := os.OpenFile(download.OutPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	defer outFile.Close()

	for download.Downloaded < download.TotalSize {

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
				decrementActiveDownload(q.ID, download.ID)
				// Block until a resume command is receivedownload.
				for {
					ctrl2 := <-controlChan
					if ctrl2 == model.ResumeCommand {
						download.Status = "downloading"
						d.UpdateDownloadInDatabase(*download)
						incrementActiveDownload(q.ID, download.ID)
						break
					}
				}
			}
		default:
			startChunk := time.Now()
			endRange := download.Downloaded + chunkSize - 1
			if download.Downloaded+chunkSize > download.TotalSize {
				endRange = download.TotalSize - 1
			}

			req, err := http.NewRequest("GET", download.URL, nil)
			if err != nil {
				// fmt.Println("Error creating GET request:", err)
				break
			}

			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", download.Downloaded, endRange))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error performing GET request:", err)
				break
			}
			defer resp.Body.Close()

			n, err := io.CopyN(outFile, resp.Body, chunkSize)
			if err != nil && err != io.EOF {
				fmt.Println("Error reading chunk:", err)
				break
			}

			elapsed := time.Since(startChunk).Seconds()
			activeCount := GetActiveDownload(download.QueueId)
			if activeCount < 1 {
				activeCount = 1
			}

			allocatedSpeed := q.MaximumBandWidth / float64(activeCount)

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
			controlChan := make(chan model.DownloadControlMessage, 1)
			ControlChannels[down.ID] = controlChan
			go d.DownloadWorker(down, ProgressChan, controlChan)
		}
	}
}
