package service

import (
	"context"

	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/model"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type DownloadService interface {
	GetDownloads(ctx context.Context) ([]entity.Download, error)
	GetDownloadById(ctx context.Context, id string) (*entity.Download, error)
	CreateDownload(ctx context.Context, dto dto.DownloadDto) (*entity.Download, error)
	DeleteDownload(ctx context.Context, id string) (*entity.Download, error)
	DownloadWorker(download *entity.Download, progressChan chan<- model.DownloadProgressMsg, controlChan <-chan model.DownloadControlMessage)
}
