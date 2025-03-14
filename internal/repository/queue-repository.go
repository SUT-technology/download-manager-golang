package repository

import (
	"context"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type QueueRepository interface {
	GetQueues(ctx context.Context) ([]entity.Queue, error)
	GetQueueById(ctx context.Context, id string) (*entity.Queue, error)
	CreateQueue(ctx context.Context, name string, savePath string, maximumDownload int, maximumBandWidth float64, activityInterval entity.TimeInterval) error
	DeleteQueue(ctx context.Context, id string) (*entity.Queue, error)
	FindAndUpdateQueue(ctx context.Context, id string, name string, savePath string, maximumDownload int, maximumBandWidth float64, activityInterval entity.TimeInterval) (*entity.Queue, error)
}
