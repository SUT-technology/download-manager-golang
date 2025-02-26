package downloadsrvc

import (
	"context"
	"fmt"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

type DownloadService struct {
	db repository.Pool
}

func NewDownloadServices(db repository.Pool) DownloadService {
	return DownloadService{db: db}
}

func (d DownloadService) GetDownloads(ctx context.Context) ([]entity.Download, error) {

	var (
		downloads []entity.Download
		err       error
	)
	queryFunc := func(r *repository.Repo) error {
		downloads, err = r.Tables.Downloads.GetDownloads(ctx)
		if err != nil {
			return fmt.Errorf("getting data from downloads: %w", err)
		}

		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	return downloads, nil
}
