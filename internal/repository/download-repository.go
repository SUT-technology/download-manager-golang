package repository

import (
	"context"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type DownloadRepository interface {
	GetDownloads(ctx context.Context) ([]entity.Download, error)
	GetDownloadById(ctx context.Context, id string) (*entity.Download, error)
	CreateDownload(ctx context.Context, url string, queueId string, fileName string , totalSize int64, outFile string) (*entity.Download, error)
	DeleteDownload(ctx context.Context, id string) (*entity.Download, error)
	UpdateDownloads(ctx context.Context, downloads []entity.Download) error
}
