package db

import (
	"context"
	"fmt"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
)

type downloadTable struct {
	pool *Pool
}

func newdownloadsTable(p *Pool) downloadTable {
	return downloadTable{pool: p}
}

func (d downloadTable) GetDownloads(ctx context.Context) ([]entity.Download, error) {
	var downloadData []entity.Download
	err := d.pool.loadData(d.pool.downloadPath, &downloadData)
	if err != nil {
		return nil, fmt.Errorf("can't load data from json: %w", err)
	}
	return downloadData, nil
}

func (d downloadTable) GetDownloadById(ctx context.Context, id string) (*entity.Download, error) {
	var downloadData []entity.Download
	err := d.pool.loadData(d.pool.downloadPath, &downloadData)
	if err != nil {
		return nil, fmt.Errorf("can't load data from json: %w", err)
	}

	var download *entity.Download
	for _, down := range downloadData {
		if down.ID == id {
			download = &down
			break
		}
	}

	return download, nil
}
