package service

import (
	"context"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type QueueService interface {
	GetQueues(ctx context.Context) ([]entity.Queue, error)
	GetQueueById(ctx context.Context, id string) (*entity.Queue, error)
	CreateQueue(ctx context.Context, dto dto.QueueDto) error
	DeleteQueue(ctx context.Context, id string) (*entity.Queue, error)
	FindAndUpdateQueue(ctx context.Context, id string, dto dto.QueueDto) error
}
