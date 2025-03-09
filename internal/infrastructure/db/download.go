package db

import (
	"context"
	"fmt"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/generator"

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
func (d downloadTable) CreateDownload(ctx context.Context, url string, queueId string, fileName string) (*entity.Download, error) {
	var downloadData []entity.Download
	err := d.pool.loadData(d.pool.downloadPath, &downloadData)
	if err != nil {
		return nil, fmt.Errorf("can't load data from json: %w", err)
	}

	ids := make([]string, len(downloadData))
	for i, download := range downloadData {
		ids[i] = download.ID
	}
	id := generator.IdGenerator(ids)

	download := &entity.Download{
		ID:       id,
		URL:      url,
		QueueId:  queueId,
		FileName: fileName,
	}

	downloadData = append(downloadData, *download)

	//TODO: implement saveData
	err = d.pool.saveData(d.pool.downloadPath, downloadData)
	if err != nil {
		return nil, fmt.Errorf("can't save data to json: %w", err)
	}

	return download, nil
}
