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

func (d downloadTable) DeleteDownload(ctx context.Context, id string) (*entity.Download, error) {
	var downloadData []entity.Download
	err := d.pool.loadData(d.pool.downloadPath, &downloadData)
	if err != nil {
		return nil, fmt.Errorf("can't load data from json: %w", err)
	}

	var indexToRemove int

	var download *entity.Download
	for i, down := range downloadData {
		if down.ID == id {
			indexToRemove = i
			download = &down
			break
		}
	}

	downloadData = append(downloadData[:indexToRemove], downloadData[indexToRemove+1:]...)

	err = d.pool.saveData(d.pool.downloadPath, downloadData)
	if err != nil {
		return nil, fmt.Errorf("can't save data to json: %w", err)
	}

	return download, nil
}

func (d downloadTable) CreateDownload(ctx context.Context, url string, queueId string, fileName string, totalSize int64, outFile string) (*entity.Download, error) {
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
		ID:           id,
		URL:          url,
		QueueId:      queueId,
		FileName:     fileName,
		TotalSize:    totalSize,
		Downloaded:   0,
		Progress:     0,
		CurrentSpeed: 0,
		Status:       "pending",
		OutPath:      outFile,
	}

	downloadData = append(downloadData, *download)

	err = d.pool.saveData(d.pool.downloadPath, downloadData)
	if err != nil {
		return nil, fmt.Errorf("can't save data to json: %w", err)
	}

	return download, nil
}

func (d downloadTable) UpdateDownloads(ctx context.Context, downloads []entity.Download) error {
	err := d.pool.saveData(d.pool.downloadPath, downloads)
	if err != nil {
		return fmt.Errorf("can't save data to json: %w", err)
	}

	return nil
}
