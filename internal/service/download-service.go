package service

import (
	"context"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type DownloadService interface {
	GetDownloads(ctx context.Context) ([]entity.Download, error)
}
