package service

import (
	"context"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type QueueService interface {
	GetQueues(ctx context.Context) ([]entity.Queue, error)
	GetQueueById(ctx context.Context, id string) (*entity.Queue, error)
}
