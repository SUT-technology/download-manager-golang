package queueHandler

import (
	"context"
	"fmt"
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
