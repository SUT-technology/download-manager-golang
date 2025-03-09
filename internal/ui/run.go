package ui

import (
	"fmt"
	"sync"

	"github.com/SUT-technology/download-manager-golang/internal/ui/model/tabs"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()
	tabs.CurrentTab.Init()
	p := tea.NewProgram(tabs.CurrentTab)
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
	return nil
}
