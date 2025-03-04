package UI

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
)

type Tab struct{
	num int
}

var CurrentTab *Tab

var tabs []Tab


func (tab Tab) Init() tea.Cmd {
	var(
		addDownloadTab = Tab{
			num:1,
		}
		downloadsListTab = Tab{
			num:2,
		}
		queuesListTab = Tab{
			num:3,
		}
	) 
	tabs=[]Tab{addDownloadTab,downloadsListTab,queuesListTab}
	CurrentTab=&downloadsListTab
	return nil
}

func (tab Tab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m:=msg.(type) {
	case tea.KeyMsg:
		if m.Type==tea.KeyShiftLeft  {
			CurrentTab = &tabs[(CurrentTab.num-1)%3]
		} else if m.Type==tea.KeyShiftRight {
			CurrentTab = &tabs[(CurrentTab.num+1)%3]
		}
	}
	return tab,nil
}

func (tab Tab) View() string {
	return fmt.Sprintf("current tab number: %v",tab.num)
}

  