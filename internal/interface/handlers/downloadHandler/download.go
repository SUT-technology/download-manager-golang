package downloadHandler

import (
	"context"
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/service"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/slogger"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

type DownloadHndlr struct {
	Services service.Services
}

func New(srvc service.Services) DownloadHndlr {
	downloadHndlr := DownloadHndlr{
		Services: srvc,
	}

	return downloadHndlr
}

func (h DownloadHndlr) GetDownloads() ([]entity.Download, error) {
	ctx := context.Background()
	slogger.Debug(ctx, "recieve request")

	downloads, err := h.Services.DownloadSrvc.GetDownloads(ctx)
	if err != nil {
		slogger.Debug(ctx, "get downloads", slogger.Err("error", err))
		return nil, fmt.Errorf("get downloads: %w", err)
	}

	return downloads, nil
}

func (h DownloadHndlr) GetDownloadById(id string) (*entity.Download, error) {
	ctx := context.Background()
	slogger.Debug(ctx, "recieve request", slog.Any("download id", id))

	download, err := h.Services.DownloadSrvc.GetDownloadById(ctx, id)
	if err != nil {
		slogger.Debug(ctx, "get download by id", slog.Any("download id", id), slogger.Err("error", err))
		return nil, fmt.Errorf("get download: %w", err)
	}

	return download, nil
}

// TODO: change outputs
func (h DownloadHndlr) CreateDownload(url string, queue *entity.Queue, fileName string) error {
	ctx := context.Background()

	var id string
	for {
		id = uuid.New().String()
		downloads, err := h.Services.DownloadSrvc.GetDownloads(ctx)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("get downloads: %w", err))
		}
		flag := false
		for _, download := range downloads {
			if download.ID == id {
				flag = true
				break
			}
		}
		if !flag {
			break
		}
	}

	download := entity.Download{
		ID:       id,
		URL:      url,
		Queue:    queue,
		FileName: fileName,
	}

	slogger.Debug(ctx, "recieve request", slog.Any("download", download))

	err := h.Services.DownloadSrvc.CreateDownload(ctx, download)
	if err != nil {
		slogger.Debug(ctx, "create download", slog.Any("download", download), slogger.Err("error", err))
		return fmt.Errorf("create download: %w", err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() error {
		defer wg.Done()
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download file: %s", resp.Status)
		}

		// Create the output file
		outFile, err := os.Create(queue.SavePath + "/" + fileName)
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
