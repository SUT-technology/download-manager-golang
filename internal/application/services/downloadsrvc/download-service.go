package downloadsrvc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/model"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

var (
	ProgressChan = make(chan model.DownloadProgressMsg)
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

	// d.StartPendingDownloads(downloads)

	return downloads, nil
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

	queryFuncQueue := func(r *repository.Repo) error {
		queue, err = r.Tables.Queues.GetQueueById(ctx, downloadDto.QueueID)
		if err != nil {
			return fmt.Errorf("Get queue by id: %w", err)
		}
		return nil
	}

	err = d.db.Query(queryFuncQueue)
	if err != nil {
		return nil, err
	}

	activeDownloads := GetActiveDownload(downloadDto.QueueID)

	if activeDownloads >= queue.MaximumDownloads {
		return nil, errors.New("downalod count for this queue is full")
	}

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

	outPath := queue.SavePath + "/" + downloadDto.FileName + "." + format
	outFile, err := os.Create(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %v", err)
	}
	outFile.Close()

	queryFunc := func(r *repository.Repo) error {

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
