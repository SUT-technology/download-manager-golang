package downloadsrvc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/model"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

var (
	// ProgressChan is where all download workers send periodic progress updates.
	ProgressChan = make(chan model.DownloadProgressMsg)
	// ControlChannels maps a download ID to its control channel.
	ControlChannels = make(map[string]chan model.DownloadControlMessage)
)

type DownloadService struct {
	db repository.Pool
}

func NewDownloadServices(db repository.Pool) DownloadService {

	downloadservice := DownloadService{db: db}

	var (
		downloads []entity.Download
		err       error
	)
	queryFunc := func(r *repository.Repo) error {
		downloads, err = r.Tables.Downloads.GetDownloads(context.Background())
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}

		return nil
	}

	err = downloadservice.db.Query(queryFunc)
	if err != nil {
		return DownloadService{}
	}

	downloadservice.StartPendingDownloads(downloads)

	return downloadservice
}

func (d DownloadService) GetDownloads(ctx context.Context) ([]entity.Download, error) {

	var (
		downloads []entity.Download
		err       error
	)
	queryFunc := func(r *repository.Repo) error {
		downloads, err = r.Tables.Downloads.GetDownloads(ctx)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}

		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	d.StartPendingDownloads(downloads)

	return downloads, nil
}

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
	download.Status = "downloading"
	d.UpdateDownloadInDatabase(*download)

	const chunkSize = 32 * 1024 // 32 KB per request
	lastBytes := download.Downloaded
	lastUpdate := time.Now()

	for download.Downloaded < download.TotalSize {
		// Check for pause/resume.
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
				// Block until a resume command is receivedownload.
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

func (d DownloadService) GetDownloadById(ctx context.Context, id string) (*entity.Download, error) {
	var (
		download *entity.Download
		err      error
	)

	queryFunc := func(r *repository.Repo) error {
		download, err = r.Tables.Downloads.GetDownloadById(ctx, id)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}

		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	return download, nil
}
func (d DownloadService) DeleteDownload(ctx context.Context, id string) (*entity.Download, error) {
	var (
		download *entity.Download
		err      error
	)

	queryFunc := func(r *repository.Repo) error {
		download, err = r.Tables.Downloads.DeleteDownload(ctx, id)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}
		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	return download, nil
}
func (d DownloadService) CreateDownload(ctx context.Context, downloadDto dto.DownloadDto) (*entity.Download, error) {

	var (
		download *entity.Download
		queue    *entity.Queue
		err      error
	)

	resp, err := http.Head(downloadDto.URL)
	if err != nil {
		return nil, fmt.Errorf("failed sending HEAD request: %v", err)
	}
	defer resp.Body.Close()

	lengthStr := resp.Header.Get("Content-Length")
	if lengthStr == "" {
		return nil, errors.New("Content-Length header not found")
	}
	totalSize, err := strconv.ParseInt(lengthStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length header: %v", err)
	}

	parsedURL, err2 := url.Parse(downloadDto.URL)
	if err2 != nil {
		return nil, err2
	}
	ext := path.Ext(parsedURL.Path)
	format := strings.TrimPrefix(ext, ".")

	if downloadDto.FileName == "" {
		downloadDto.FileName = path.Base(parsedURL.Path)
	}

	if err != nil {
		return nil, err
	}

	queryFunc := func(r *repository.Repo) error {
		queue, err = r.Tables.Queues.GetQueueById(ctx, downloadDto.QueueID)
		if err != nil {
			return fmt.Errorf("Get queue by id: %w", err)
		}

		outPath := queue.SavePath + "/" + downloadDto.FileName + "." + format
		outFile, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		outFile.Close()
		download, err = r.Tables.Downloads.CreateDownload(ctx, downloadDto.URL, downloadDto.QueueID, downloadDto.FileName, totalSize, outPath)
		if err != nil {
			return fmt.Errorf("creating download: %w", err)
		}
		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	controlChan := make(chan model.DownloadControlMessage)
	ControlChannels[download.ID] = controlChan
	go d.DownloadWorker(download, ProgressChan, controlChan)

	return download, nil
}
