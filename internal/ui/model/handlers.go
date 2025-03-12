package model

import (
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers/downloadHandler"
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers/queueHandler"
)

type Handlers struct {
	DownloadHandler *downloadHandler.DownloadHndlr
	QueueHandler    *queueHandler.QueueHndlr
}
