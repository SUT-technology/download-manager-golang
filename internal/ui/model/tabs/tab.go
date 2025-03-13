package tabs

import (
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/ui/model"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type (
	errMsg error
)

type Tab struct {
	num int
	TAB interface{}
}

type AddDownloadTab struct {
	downloads []entity.Download
	queues    []entity.Queue
	adm       addDownloadModel
	err       error
}

type DownloadsTab struct {
	downloads []entity.Download
}

type queuesTableTab struct {
	queues []entity.Queue
}

type addDownloadModel struct {
	url             string
	urlInput        textinput.Model
	selectedQueueId string
	fileName        string
	fileNameInput   textinput.Model
	cursorIndex     int
	finished        bool
}

var hndlr model.Handlers

func InitiateAddDownloadTab(Hndlr *model.Handlers) AddDownloadTab {

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

	return AddDownloadTab{
		downloads: downloads,
		queues:    queues,
		adm: addDownloadModel{
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

func InitiateDownloadsTab(Hndlr *model.Handlers) DownloadsTab {
	hndlr := *Hndlr
	downloads, err := hndlr.DownloadHandler.GetDownloads()
	if err != nil {
		panic(err)
	}
	downloadsTab := DownloadsTab{
		downloads: downloads,
	}
	return downloadsTab
}

func (tab Tab) Init() tea.Cmd {
	return textinput.Blink
}

func (tab Tab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if tab.num == 2 {
		var addDownloadTab = tab.TAB.(AddDownloadTab)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return tab, tea.Quit
			}

		// We handle errors just like any other message
		case errMsg:
			addDownloadTab.err = msg
			return addDownloadTab, nil
		}

		if addDownloadTab.adm.finished {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if key := msg.String(); strings.ToLower(key)[0] == 'y' {
					addDownloadTab.adm.finished = false
					return InitiateAddDownloadTab(&hndlr), nil
				} else if strings.ToLower(key)[0] == 'n' {
					return addDownloadTab, tea.Quit
				}
			}
		} else if addDownloadTab.adm.url == "" {

			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					addDownloadTab.adm.url = addDownloadTab.adm.urlInput.Value()
				}
				addDownloadTab.adm.urlInput, cmd = addDownloadTab.adm.urlInput.Update(msg)
			}

		} else if addDownloadTab.adm.selectedQueueId == "" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyUp:
					if addDownloadTab.adm.cursorIndex == 0 {
						addDownloadTab.adm.cursorIndex = len(addDownloadTab.queues) - 1
					} else {
						addDownloadTab.adm.cursorIndex--
					}
				case tea.KeyDown:
					addDownloadTab.adm.cursorIndex = (addDownloadTab.adm.cursorIndex + 1) % len(addDownloadTab.queues)
				case tea.KeyEnter:
					addDownloadTab.adm.selectedQueueId = addDownloadTab.queues[addDownloadTab.adm.cursorIndex].ID
				}
			}

		} else if addDownloadTab.adm.fileName == "" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					addDownloadTab.adm.fileName = addDownloadTab.adm.fileNameInput.Value()

					err := CreateDownload(addDownloadTab.adm.url, addDownloadTab.adm.selectedQueueId, addDownloadTab.adm.fileName)
					if err != nil {
						addDownloadTab.err = err
						//return addDownloadTab, nil
					}

					addDownloadTab.adm.finished = true
					return addDownloadTab, nil
				}
				addDownloadTab.adm.fileNameInput, cmd = addDownloadTab.adm.fileNameInput.Update(msg)
			}
		}

	} else if tab.num == 1 {

	} else if tab.num == 3 {

	}
	return tab, cmd
}

func (tab Tab) View() string {
	var view string
	if tab.num == 1 {
		var addDownloadTab = tab.TAB.(AddDownloadTab)
		if addDownloadTab.adm.finished {
			view = "Download added successfully!\n\nDo you want to continue? (y/n)"
		} else if addDownloadTab.adm.url == "" {
			view = fmt.Sprintf(
				"\nEnter the url here:\n\n%s\n\n%s",
				addDownloadTab.adm.urlInput.View(),
				"(ctrl+c to quit)",
			) + "\n"
		} else if addDownloadTab.adm.selectedQueueId == "" {
			view = "\nSelect a queue:\n\n"
			for i, queue := range addDownloadTab.queues {
				cursor := " "
				if i == addDownloadTab.adm.cursorIndex {
					cursor = ">"
				}
				view += fmt.Sprintf("%s %s\n", cursor, queue.Name)
			}
		} else if addDownloadTab.adm.fileName == "" {
			view = fmt.Sprintf(
				"\nEnter the file name here (optional):\n\n%s\n\n%s",
				addDownloadTab.adm.fileNameInput.View(),
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
