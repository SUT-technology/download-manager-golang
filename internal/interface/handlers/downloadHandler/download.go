package downloadHandler

import (
	"context"
	"fmt"

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
