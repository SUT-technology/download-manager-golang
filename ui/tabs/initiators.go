package tabs

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func InitiateAddDownloadTab(Hndlr *Handlers) Tab {

	hndlr = *Hndlr

	uInp := textinput.New()
	uInp.Placeholder = "url"
	uInp.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	uInp.Focus()
	uInp.Width = 100

	downloads, err := hndlr.DownloadHandler.GetDownloads()
	//Todo fix err
	if err != nil {
		panic(err)
	}

	queues, err := hndlr.QueueHandler.GetQueues()
	//Todo fix err
	if err != nil {
		panic(err)
	}

	fInp := textinput.New()
	fInp.Placeholder = "Pikachu"
	fInp.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	fInp.Focus()
	fInp.Width = 100

	return Tab{
		num: 1,
		TAB: AddDownloadTab{
			downloads:       downloads,
			queues:          queues,
			url:             "",
			urlInput:        uInp,
			selectedQueueId: "",
			fileName:        "",
			fileNameInput:   fInp,
			cursorIndex:     0,
			finished:        false,
		},
	}
}

func InitiateDownloadsTab(Hndlr *Handlers) Tab {
	hndlr = *Hndlr
	downloads, err := hndlr.DownloadHandler.GetDownloads()
	if err != nil {
		panic(err)
	}
	return Tab{
		num: 2,
		TAB: DownloadsTab{
			downloads:    downloads,
			cursorIndex:  0,
			deleteAction: false,
			message: "",
		},
	}
}

func InitiateQueuesTab(Hndlr *Handlers) Tab {
	hndlr = *Hndlr

	nameInput := textinput.New()
	nameInput.Placeholder = ""
	nameInput.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	nameInput.Focus()
	nameInput.Width = 100

	savePathInput := textinput.New()
	savePathInput.Placeholder = ""
	savePathInput.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	savePathInput.Focus()
	savePathInput.Width = 100

	maximumDownloadsInput := textinput.New()
	maximumDownloadsInput.Placeholder = ""
	maximumDownloadsInput.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	maximumDownloadsInput.Focus()
	maximumDownloadsInput.Width = 100

	maximumBandWidthInput := textinput.New()
	maximumBandWidthInput.Placeholder = ""
	maximumBandWidthInput.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	maximumBandWidthInput.Focus()
	maximumBandWidthInput.Width = 100

	startTimeInput := textinput.New()
	startTimeInput.Placeholder = ""
	startTimeInput.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	startTimeInput.Focus()
	startTimeInput.Width = 100

	endTimeInput := textinput.New()
	endTimeInput.Placeholder = ""
	endTimeInput.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	endTimeInput.Focus()
	endTimeInput.Width = 100

	queues, err := hndlr.QueueHandler.GetQueues()
	if err != nil {
		panic(err)
	}
	return Tab{
		num: 3,
		TAB: QueuesTab{
			queues:                queues,
			cursorIndex:           0,
			action:                "list",
			nameInput:             nameInput,
			savePathInput:         savePathInput,
			maximumDownloadsInput: maximumDownloadsInput,
			maximumBandWidthInput: maximumBandWidthInput,
			startTimeInput:        startTimeInput,
			endTimeInput:          endTimeInput,
		},
	}
}
