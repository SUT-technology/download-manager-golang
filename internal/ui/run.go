package ui

import (
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers"
	"sync"

	//"github.com/SUT-technology/download-manager-golang/internal/interface/handlers/downloadHandler"
	//"github.com/SUT-technology/download-manager-golang/internal/interface/handlers/queueHandler"
	"github.com/SUT-technology/download-manager-golang/internal/ui/model"
	"github.com/SUT-technology/download-manager-golang/internal/ui/model/tabs"
	tea "github.com/charmbracelet/bubbletea"
)

var Hndlr model.Handlers

func Run(wg *sync.WaitGroup, srcv *handlers.HandlerSrcv) error {
	//wg.Add(1)
	defer wg.Done()

	Hndlr = model.Handlers{
		DownloadHandler: &srcv.DownloadHndlr,
		QueueHandler:    &srcv.QueueHndlr,
	}

	p := tea.NewProgram(tabs.InitiateDownloadsTab(&Hndlr))
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
	return nil
}
