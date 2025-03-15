package service

import (
	"context"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type DownloadService interface {
	GetDownloads(ctx context.Context) ([]entity.Download, error)
	GetDownloadById(ctx context.Context, id string) (*entity.Download, error)
	CreateDownload(ctx context.Context, dto dto.DownloadDto) error
	DeleteDownload(ctx context.Context, id string) (*entity.Download, error)
}
