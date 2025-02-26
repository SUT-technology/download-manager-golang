package repository

import (
	"context"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type DownloadRepository interface {
	GetDownloads(ctx context.Context) ([]entity.Download , error)
}