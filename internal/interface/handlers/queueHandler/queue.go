package queueHandler

import (
	"context"
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"log/slog"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/service"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/slogger"
)

type QueueHndlr struct {
	Services service.Services
}

func New(srvc service.Services) QueueHndlr {
	queueHndlr := QueueHndlr{
		Services: srvc,
	}

	return queueHndlr
}

func (h QueueHndlr) GetQueues() ([]entity.Queue, error) {
	ctx := context.Background()
	slogger.Debug(ctx, "recieve request")

	queues, err := h.Services.QueueSrvc.GetQueues(ctx)
	if err != nil {
		slogger.Debug(ctx, "get queues", slogger.Err("error", err))
		return nil, fmt.Errorf("get queues: %w", err)
	}

	return queues, nil
}

func (h QueueHndlr) GetQueueById(id string) (*entity.Queue, error) {
	ctx := context.Background()
	slogger.Debug(ctx, "recieve request", slog.Any("queue id", id))

	queue, err := h.Services.QueueSrvc.GetQueueById(ctx, id)
	if err != nil {
		slogger.Debug(ctx, "get queue by id", slog.Any("queue id", id), slogger.Err("error", err))
		return nil, fmt.Errorf("get queue: %w", err)
	}

	return queue, nil
}

func (h QueueHndlr) CreateQueue(queueDto dto.QueueDto) error {
	ctx := context.Background()

	slogger.Debug(ctx, "recieve request", slog.Any("queueDto", queueDto))

	err := h.Services.QueueSrvc.CreateQueue(ctx, queueDto.Name, queueDto.SavePath, queueDto.MaximumDownloads, queueDto.MaximumBandWidth, queueDto.ActivityInterval)
	if err != nil {
		slogger.Debug(ctx, "create queue", slog.Any("queueDto", queueDto), slogger.Err("error", err))
		return fmt.Errorf("create queue: %w", err)
	}

	return nil
}
