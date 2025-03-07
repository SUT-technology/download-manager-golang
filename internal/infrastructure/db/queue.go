package db

import (
	"context"
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type queueTable struct {
	pool *Pool
}

func newqueuesTable(p *Pool) queueTable {
	return queueTable{pool: p}
}

func (d queueTable) GetQueues(ctx context.Context) ([]entity.Queue, error) {
	var queueData []entity.Queue
	err := d.pool.loadData(d.pool.queuePath, &queueData)
	if err != nil {
		return nil, fmt.Errorf("can't load data from json: %w", err)
	}
	return queueData, nil
}

func (d queueTable) GetQueueById(ctx context.Context, id string) (*entity.Queue, error) {
	var queueData []entity.Queue
	err := d.pool.loadData(d.pool.queuePath, &queueData)
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
