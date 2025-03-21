package queuesrvc

import (
	"context"
	"fmt"
	"os"

	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

type QueueService struct {
	db repository.Pool
}

func (q QueueService) DeleteQueue(ctx context.Context, id string) (*entity.Queue, error) {
	var (
		queue *entity.Queue
		err   error
	)

	queryFunc := func(r *repository.Repo) error {
		queue, err = r.Tables.Queues.DeleteQueue(ctx, id)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}
		return nil
	}

	err = q.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	return queue, nil
}

func NewQueueServices(db repository.Pool) QueueService {
	return QueueService{db: db}
}

func (q QueueService) GetQueues(ctx context.Context) ([]entity.Queue, error) {

	var (
		queues []entity.Queue
		err    error
	)
	queryFunc := func(r *repository.Repo) error {
		queues, err = r.Tables.Queues.GetQueues(ctx)
		if err != nil {
			return fmt.Errorf("getting data from queues: %w", err)
		}

		return nil
	}

	err = q.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	return queues, nil
}

func (q QueueService) GetQueueById(ctx context.Context, id string) (*entity.Queue, error) {
	var (
		queue *entity.Queue
		err   error
	)

	queryFunc := func(r *repository.Repo) error {
		queue, err = r.Tables.Queues.GetQueueById(ctx, id)
		if err != nil {
			return fmt.Errorf("getting data from queues: %w", err)
		}

		return nil
	}

	err = q.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	return queue, nil
}

func (q QueueService) CreateQueue(ctx context.Context, dto dto.QueueDto) error {
	var (
		err error
	)

	queryFunc := func(r *repository.Repo) error {
		err = r.Tables.Queues.CreateQueue(ctx, dto.Name, dto.SavePath, dto.MaximumDownloads, dto.MaximumBandWidth, dto.ActivityInterval)
		if err != nil {
			return fmt.Errorf("creating queue: %w", err)
		}

		return nil
	}

	err = q.db.Query(queryFunc)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dto.SavePath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directories: %w", err)
	}

	return nil
}

func (q QueueService) FindAndUpdateQueue(ctx context.Context, id string, dto dto.QueueDto) error {
	var err error

	queryFunc := func(r *repository.Repo) error {
		err = r.Tables.Queues.FindAndUpdateQueue(ctx, id, dto.Name, dto.SavePath, dto.MaximumDownloads, dto.MaximumBandWidth, dto.ActivityInterval)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}
		return nil
	}

	err = q.db.Query(queryFunc)
	if err != nil {
		return err
	}

	return nil
}
