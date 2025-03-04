package UI

import(
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/SUT-technology/download-manager-golang/internal/UI/model/tabs"
) 


func Run() error {
	UI.CurrentTab.Init()
	p := tea.NewProgram(UI.CurrentTab)  
    if err := p.Start(); err != nil {  
        fmt.Printf("Error starting program: %v\n", err)  
    }
	return nil
}