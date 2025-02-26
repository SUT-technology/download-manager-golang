package handlers

import (
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers/downloadHandler"
	"github.com/SUT-technology/download-manager-golang/internal/service"
)

type HandlerSrcv struct {
	DownloadHndlr downloadHandler.DownloadHndlr
}

func New(srvc service.Services) HandlerSrcv {
	return HandlerSrcv{
		DownloadHndlr: downloadHandler.New(srvc),
	}
}
