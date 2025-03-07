package handlers

import (
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers/downloadHandler"
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers/queueHandler"
	"github.com/SUT-technology/download-manager-golang/internal/service"
)

type HandlerSrcv struct {
	DownloadHndlr downloadHandler.DownloadHndlr
	QueueHndlr    queueHandler.QueueHndlr
}

func New(srvc service.Services) HandlerSrcv {
	return HandlerSrcv{
		DownloadHndlr: downloadHandler.New(srvc),
		QueueHndlr:    queueHandler.New(srvc),
	}
}
