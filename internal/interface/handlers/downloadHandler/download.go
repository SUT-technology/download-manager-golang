package downloadHandler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/domain/model"
	"github.com/SUT-technology/download-manager-golang/internal/service"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/slogger"
	//"io"
	//"net/http"
	//"os"
	//"sync"
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

func (h DownloadHndlr) DeleteDownload(id string) (*entity.Download, error) {
	ctx := context.Background()
	slogger.Debug(ctx, "recieve request", slog.Any("download id", id))

	download, err := h.Services.DownloadSrvc.DeleteDownload(ctx, id)
	if err != nil {
		slogger.Debug(ctx, "delete download by id", slog.Any("download id", id), slogger.Err("error", err))
		return nil, fmt.Errorf("get download: %w", err)
	}

	return download, nil
}

func (h DownloadHndlr) CreateDownload(downloadDto dto.DownloadDto) (*entity.Download, error) {
	ctx := context.Background()

	slogger.Debug(ctx, "recieve request", slog.Any("download", downloadDto))

	download, err := h.Services.DownloadSrvc.CreateDownload(ctx, downloadDto)
	if err != nil {
		slogger.Debug(ctx, "create download", slog.Any("download", downloadDto), slogger.Err("error", err))
		return nil, fmt.Errorf("create download: %w", err)
	}

	return download, nil
}

func (h DownloadHndlr) DownloadWorker(download *entity.Download, progressChan chan<- model.DownloadProgressMsg, controlChan <-chan model.DownloadControlMessage) {
	ctx := context.Background()

	slogger.Debug(ctx, "recieve request", slog.Any("download", download))

	h.Services.DownloadSrvc.DownloadWorker(download, progressChan, controlChan)

	return
}
