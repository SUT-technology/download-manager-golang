package services

import (
	"github.com/SUT-technology/download-manager-golang/internal/application/services/downloadsrvc"
	"github.com/SUT-technology/download-manager-golang/internal/application/services/queuesrvc"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
	"github.com/SUT-technology/download-manager-golang/internal/service"
)

func New(db repository.Pool) service.Services {
	return service.Services{
		DownloadSrvc: downloadsrvc.NewDownloadServices(db),
		QueueSrvc:    queuesrvc.NewQueueServices(db),
	}
}
