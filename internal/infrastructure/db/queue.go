package db

import (
	"context"
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/generator"
)

type queueTable struct {
	pool *Pool
}

func newqueuesTable(p *Pool) queueTable {
	return queueTable{pool: p}
}

func (q queueTable) GetQueues(ctx context.Context) ([]entity.Queue, error) {
	var queueData []entity.Queue
	err := q.pool.loadData(q.pool.queuePath, &queueData)
	if err != nil {
		return nil, fmt.Errorf("can't load data from json: %w", err)
	}
	return queueData, nil
}

func (q queueTable) GetQueueById(ctx context.Context, id string) (*entity.Queue, error) {
	var queueData []entity.Queue
	err := q.pool.loadData(q.pool.queuePath, &queueData)
	if err != nil {
		return nil, fmt.Errorf("can't load data from json: %w", err)
	}

	var queue *entity.Queue
	for _, q := range queueData {
		if q.ID == id {
			queue = &q
			break
		}
	}

	return queue, nil
}

func (q queueTable) CreateQueue(ctx context.Context, name string, savePath string, maximumDownload int, maximumBandWidth float64, activityInterval entity.TimeInterval) error {
	var queueData []entity.Queue
	err := q.pool.loadData(q.pool.queuePath, &queueData)
	if err != nil {
		return fmt.Errorf("can't load data from json: %w", err)
	}

	ids := make([]string, len(queueData))
	for i, q := range queueData {
		ids[i] = q.ID
	}
	id := generator.IdGenerator(ids)

	queue := entity.Queue{
		ID:               id,
		Name:             name,
		SavePath:         savePath,
		MaximumDownloads: maximumDownload,
		MaximumBandWidth: maximumBandWidth,
		ActivityInterval: activityInterval,
	}

	queueData = append(queueData, queue)

	err = q.pool.saveData(q.pool.queuePath, queueData)
	if err != nil {
		return fmt.Errorf("can't save data to json: %w", err)
	}

	return nil
}
