package downloadHandler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/service"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/slogger"
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
