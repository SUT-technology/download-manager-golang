package tabs

import (
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/ui/model"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"runtime"
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
	downloads       []entity.Download
	queues          []entity.Queue
	url             string
	urlInput        textinput.Model
	selectedQueueId string
	fileName        string
	fileNameInput   textinput.Model
	cursorIndex     int
	finished        bool
	err             error
}

func InitiateAddDownloadTab(Hndlr *model.Handlers) Tab {

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

type DownloadsTab struct {
	downloads    []entity.Download
	cursorIndex  int
	deleteAction bool
	message      string
}

func InitiateDownloadsTab(Hndlr *model.Handlers) Tab {
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
		},
	}
}

type QueuesTab struct {
	queues           []entity.Queue
	cursorIndex      int
	deleteAction     bool
	editAction       bool
	maxSpeed         textinput.Model
	savePath         textinput.Model
	maximumDownloads textinput.Model
	bandWidth        textinput.Model
	message          string
}

func InitiateQueuesTab(Hndlr *model.Handlers) Tab {
	hndlr = *Hndlr

	maxSpeed := textinput.New()
	maxSpeed.Placeholder = ""
	maxSpeed.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	maxSpeed.Focus()
	maxSpeed.Width = 100

	savePath := textinput.New()
	savePath.Placeholder = ""
	savePath.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	savePath.Focus()
	savePath.Width = 100

	maximumDownloads := textinput.New()
	maximumDownloads.Placeholder = ""
	maximumDownloads.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	maximumDownloads.Focus()
	maximumDownloads.Width = 100

	bandWidth := textinput.New()
	bandWidth.Placeholder = ""
	bandWidth.TextStyle = lipgloss.NewStyle().Blink(true).Foreground(lipgloss.Color("205"))
	bandWidth.Focus()
	bandWidth.Width = 100

	queues, err := hndlr.QueueHandler.GetQueues()
	if err != nil {
		panic(err)
	}
	return Tab{
		num: 2,
		TAB: QueuesTab{
			queues:           queues,
			cursorIndex:      0,
			deleteAction:     false,
			editAction:       false,
			maxSpeed:         maxSpeed,
			savePath:         savePath,
			maximumDownloads: maximumDownloads,
			bandWidth:        bandWidth,
		},
	}
}

var hndlr model.Handlers

func (tab Tab) Init() tea.Cmd {
	return textinput.Blink
}

func (tab Tab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyShiftLeft:
			if tab.num == 1 {
				ClearScreen()
				return InitiateDownloadsTab(&hndlr), cmd
			} else if tab.num == 2 {
				ClearScreen()
				return InitiateAddDownloadTab(&hndlr), cmd
			}

		case tea.KeyShiftRight:
		}
	}
	if tab.num == 1 {
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
			return tab, nil
		}

		if addDownloadTab.finished {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if key := msg.String(); strings.ToLower(key)[0] == 'y' {
					addDownloadTab.finished = false
					return InitiateAddDownloadTab(&hndlr), nil
				} else if strings.ToLower(key)[0] == 'n' {
					return tab, tea.Quit
				}
			}
		} else if addDownloadTab.url == "" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					addDownloadTab.url = addDownloadTab.urlInput.Value()
				}
				addDownloadTab.urlInput, cmd = addDownloadTab.urlInput.Update(msg)
			}
		} else if addDownloadTab.selectedQueueId == "" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyUp:
					if addDownloadTab.cursorIndex == 0 {
						addDownloadTab.cursorIndex = len(addDownloadTab.queues) - 1
					} else {
						addDownloadTab.cursorIndex--
					}
				case tea.KeyDown:
					addDownloadTab.cursorIndex = (addDownloadTab.cursorIndex + 1) % len(addDownloadTab.queues)
				case tea.KeyEnter:
					addDownloadTab.selectedQueueId = addDownloadTab.queues[addDownloadTab.cursorIndex].ID
				}
			}
		} else if addDownloadTab.fileName == "" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					addDownloadTab.fileName = addDownloadTab.fileNameInput.Value()
					err := CreateDownload(addDownloadTab.url, addDownloadTab.selectedQueueId, addDownloadTab.fileName)
					if err != nil {
						addDownloadTab.err = err
					}
					addDownloadTab.finished = true
					tab.TAB = addDownloadTab
					return tab, nil
				}
				addDownloadTab.fileNameInput, cmd = addDownloadTab.fileNameInput.Update(msg)
			}
		}
		tab.TAB = addDownloadTab
		return tab, cmd
	} else if tab.num == 2 {
		downloadsTab := tab.TAB.(DownloadsTab)
		if !downloadsTab.deleteAction {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyUp:
					if downloadsTab.cursorIndex == 0 {
						downloadsTab.cursorIndex = len(downloadsTab.downloads) - 1
					} else {
						downloadsTab.cursorIndex--
					}
				case tea.KeyDown:
					downloadsTab.cursorIndex = (downloadsTab.cursorIndex + 1) % len(downloadsTab.downloads)
				case tea.KeyCtrlA:

				//TODO: pause/resume

				case tea.KeyCtrlS:

				//TODO: retry

				case tea.KeyCtrlD:
					download, err := hndlr.DownloadHandler.DeleteDownload(downloadsTab.downloads[downloadsTab.cursorIndex].ID)
					if err != nil {
						panic(err)
					}
					downloadsTab.message = fmt.Sprintf("%s deleted successfully", download.FileName)
					downloadsTab.deleteAction = true
				}
			}
		} else {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if key := msg.String(); strings.ToLower(key)[0] == 'y' {
					downloadsTab.deleteAction = false
					return InitiateDownloadsTab(&hndlr), nil
				} else if strings.ToLower(key)[0] == 'n' {
					return tab, tea.Quit
				}
			}
		}
		tab.TAB = downloadsTab
		return tab, cmd
	} else if tab.num == 3 {
		queuesTab := tab.TAB.(QueuesTab)
		if !queuesTab.deleteAction && !queuesTab.editAction {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyUp:
					if queuesTab.cursorIndex == 0 {
						queuesTab.cursorIndex = len(queuesTab.queues) - 1
					} else {
						queuesTab.cursorIndex--
					}
				case tea.KeyDown:
					queuesTab.cursorIndex = (queuesTab.cursorIndex + 1) % len(queuesTab.queues)
				}
			}
		}
		tab.TAB = queuesTab
		return tab, cmd
	}
	return tab, cmd
}

func (tab Tab) View() string {
	var view string
	if tab.num == 1 {
		view = "                   ------------------------ Add Download Tab ------------------------"
		var addDownloadTab = tab.TAB.(AddDownloadTab)
		if addDownloadTab.finished {
			view += "\nDownload added successfully!\n\nDo you want to continue? (y => yes) (n => no)"
		} else if addDownloadTab.url == "" {
			view += fmt.Sprintf(
				"\nEnter the url here:\n\n%s\n\n%s",
				addDownloadTab.urlInput.View(),
				"(ctrl+c to quit)",
			) + "\n"
		} else if addDownloadTab.selectedQueueId == "" {
			view += "\nSelect a queue:\n\n"
			for i, queue := range addDownloadTab.queues {
				cursor := "  "
				if i == addDownloadTab.cursorIndex {
					cursor = "> "
				}
				view += fmt.Sprintf("%s%s\n", cursor, queue.Name)
			}
		} else if addDownloadTab.fileName == "" {
			view += fmt.Sprintf(
				"\nEnter the file name here (optional):\n\n%s\n\n%s",
				addDownloadTab.fileNameInput.View(),
				"(ctrl+c => quit)",
			) + "\n"
		}
	} else if tab.num == 2 {
		var downloadsTab = tab.TAB.(DownloadsTab)
		view = "                   ------------------------ Downloads Tab ------------------------"
		if !downloadsTab.deleteAction {
			view = fmt.Sprintf("%v\nSelect a queue:", view)
			for i, download := range downloadsTab.downloads {
				queue, err := hndlr.QueueHandler.GetQueueById(download.QueueId)
				if err != nil {
					panic(err)
				}
				cursor := "  "
				if i == downloadsTab.cursorIndex {
					cursor = "> "
				}
				view = fmt.Sprintf("%v\n%vURL: %v    Queue: %v    Status: %v    Speed: %v    Progress: %v", view, cursor, download.URL, queue.Name, download.Status, download.CurrentSpeed, download.Progress)
			}
			view = fmt.Sprintf("%v\n\n(ctrl+a => pause/resume) (ctrl+s => retry) (ctrl+d => delete)", view)
		}
		if downloadsTab.deleteAction {
			view = fmt.Sprintf("%v\n%v\n%v", view, downloadsTab.message, "do you want to continue? (y => yes) (n => no)")
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

func ClearScreen() {
	var cmd *exec.Cmd
	// Check the OS type and run the corresponding clear command
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // for Windows
	default:
		cmd = exec.Command("clear") // for Unix-like systems (Linux, macOS)
	}
	// Run the clear command
	cmd.Stdout = os.Stdout
	cmd.Run()
}
