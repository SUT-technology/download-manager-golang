package downloadsrvc

import (
	"context"
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/application/services/queuesrvc"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

type DownloadService struct {
	db repository.Pool
}

func NewDownloadServices(db repository.Pool) DownloadService {
	return DownloadService{db: db}
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

func (d DownloadService) CreateDownload(ctx context.Context, Url string, queueId string, fileName string) error {
	var (
		download *entity.Download
		err      error
	)

	queryFunc := func(r *repository.Repo) error {
		download, err = r.Tables.Downloads.CreateDownload(ctx, Url, queueId, fileName)
		if err != nil {
			return fmt.Errorf("creating download: %w", err)
		}

		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return err
	}

	queueSrvc := queuesrvc.NewQueueServices(d.db)

	queue, err := queueSrvc.GetQueueById(ctx, download.QueueId)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() error {
		defer wg.Done()
		resp, err := http.Get(download.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download file: %s", resp.Status)
		}

		parsedURL, err2 := url.Parse(download.URL)
		if err2 != nil {
			return err2
		}
		ext := path.Ext(parsedURL.Path)        // Get extension (e.g., .mp4)
		format := strings.TrimPrefix(ext, ".") // Remove leading dot

		if download.FileName == "" {
			download.FileName = path.Base(parsedURL.Path)
		}

		// Create the output file
		outFile, err := os.Create(queue.SavePath + "/" + download.FileName + "." + format)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Copy the response body to the file
		_, err = io.Copy(outFile, resp.Body)
		return err
	}()

	wg.Wait()

	return nil
}
