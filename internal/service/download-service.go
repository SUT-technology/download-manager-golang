package service

import (
	"context"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type DownloadService interface {
	GetDownloads(ctx context.Context) ([]entity.Download, error)
	GetDownloadById(ctx context.Context, id string) (*entity.Download, error)
	CreateDownload(ctx context.Context, download entity.Download) error
}
