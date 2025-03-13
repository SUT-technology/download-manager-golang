package tabs

import (
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/charmbracelet/lipgloss"
	"strings"

	//"github.com/SUT-technology/download-manager-golang/internal/ui"

	//"github.com/SUT-technology/download-manager-golang/internal/ui"
	"github.com/SUT-technology/download-manager-golang/internal/ui/model"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type Tab struct {
	num       int
	downloads []entity.Download
	queues    []entity.Queue
	adm       addDownloadModel
	err       error
}

type addDownloadModel struct {
	url      string
	urlInput textinput.Model
	//queues        []string
	selectedQueueId string
	fileName        string
	fileNameInput   textinput.Model
	cursorIndex     int
	finished        bool
}

var CurrentTab *Tab

var hndlr model.Handlers

var tabs []Tab

func InitialModel(Hndlr *model.Handlers) Tab {
	hndlr = *Hndlr

	uInp := textinput.New()
	uInp.Placeholder = "url"
	//uInp.PlaceholderStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("180"))
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

	//queuesString := make([]string, len(queues))
	//for i, v := range queues {
	//	queuesString[i] = v.Name
	//}

	fInp := textinput.New()
	fInp.Placeholder = "Pikachu"
	//fInp.PlaceholderStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("180"))
	fInp.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	fInp.Focus()
	fInp.Width = 100

	return Tab{
		num:       1,
		downloads: downloads,
		queues:    queues,
		adm: addDownloadModel{
			url:      "",
			urlInput: uInp,
			//queues:        queuesString,
			selectedQueueId: "",
			fileName:        "",
			fileNameInput:   fInp,
			cursorIndex:     0,
			finished:        false,
		},
	}
}

func (tab Tab) Init() tea.Cmd {
	return textinput.Blink
}

func (tab Tab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if tab.num == 1 {

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return tab, tea.Quit
			}

		// We handle errors just like any other message
		case errMsg:
			tab.err = msg
			return tab, nil
		}

		if tab.adm.finished {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if key := msg.String(); strings.ToLower(key)[0] == 'y' {
					tab.adm.finished = false
					return InitialModel(&hndlr), nil
				} else if strings.ToLower(key)[0] == 'n' {
					return tab, tea.Quit
				}
			}
		} else if tab.adm.url == "" {

			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					tab.adm.url = tab.adm.urlInput.Value()
				}
				tab.adm.urlInput, cmd = tab.adm.urlInput.Update(msg)
			}

		} else if tab.adm.selectedQueueId == "" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyUp:
					if tab.adm.cursorIndex == 0 {
						tab.adm.cursorIndex = len(tab.queues) - 1
					} else {
						tab.adm.cursorIndex--
					}
				case tea.KeyDown:
					tab.adm.cursorIndex = (tab.adm.cursorIndex + 1) % len(tab.queues)
				case tea.KeyEnter:
					tab.adm.selectedQueueId = tab.queues[tab.adm.cursorIndex].ID
				}
			}

		} else if tab.adm.fileName == "" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					tab.adm.fileName = tab.adm.fileNameInput.Value()

					err := CreateDownload(tab.adm.url, tab.adm.selectedQueueId, tab.adm.fileName)
					if err != nil {
						tab.err = err
						//return tab, nil
					}

					tab.adm.finished = true
					return tab, nil
				}
				tab.adm.fileNameInput, cmd = tab.adm.fileNameInput.Update(msg)
			}
		}

	}
	return tab, cmd
}

func (tab Tab) View() string {
	var view string
	if tab.num == 1 {
		if tab.adm.finished {
			view = "Download added successfully!\n\nDo you want to continue? (y/n)"
		} else if tab.adm.url == "" {
			view = fmt.Sprintf(
				"\nEnter the url here:\n\n%s\n\n%s",
				tab.adm.urlInput.View(),
				"(ctrl+c to quit)",
			) + "\n"
		} else if tab.adm.selectedQueueId == "" {
			//view = fmt.Sprintf(
			//	"Enter the url here\n\n%s\n\n%s",
			//	tab.adm.urlInput.View(),
			//	"(esc to quit)",
			//) + "\n\n"
			view = "\nSelect a queue:\n\n"
			for i, queue := range tab.queues {
				cursor := " "
				if i == tab.adm.cursorIndex {
					cursor = ">"
				}
				view += fmt.Sprintf("%s %s\n", cursor, queue.Name)
			}
		} else if tab.adm.fileName == "" {
			view = fmt.Sprintf(
				"\nEnter the file name here (optional):\n\n%s\n\n%s",
				tab.adm.fileNameInput.View(),
				"(ctrl+c to quit)",
			) + "\n"
		}
	}

	return view
}

func CreateDownload(url, queueId, fileName string) error {
	return hndlr.DownloadHandler.CreateDownload(dto.DownloadDto{
		URL:      url,
		QueueID:  queueId,
		FileName: fileName,
	})
}
