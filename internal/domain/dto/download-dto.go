package dto

import "github.com/SUT-technology/download-manager-golang/internal/domain/entity"

type DownloadDto struct {
	URL      string
	Queue    *entity.Queue
	FileName string
}
