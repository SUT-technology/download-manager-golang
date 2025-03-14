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
	"strconv"
	"strings"
	"time"
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
	queues                []entity.Queue
	cursorIndex           int
	action                string
	nameInput             textinput.Model
	savePathInput         textinput.Model
	maximumDownloadsInput textinput.Model
	maximumBandWidthInput textinput.Model
	startTimeInput        textinput.Model
	endTimeInput          textinput.Model
	name                  string
	savePath              string
	maximumDownloads      int
	maximumBandWidth      float64
	startTime             string
	endTime               string
	message               string
	err                   error
}

func InitiateQueuesTab(Hndlr *model.Handlers) Tab {
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

var hndlr model.Handlers

func (tab Tab) Init() tea.Cmd {
	return textinput.Blink
}

func (tab Tab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch message := msg.(type) {

	// general key listeners: shifting between tabs and quit
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyShiftLeft:
			if tab.num == 3 {
				ClearScreen()
				return InitiateDownloadsTab(&hndlr), cmd
			} else if tab.num == 1 {
				ClearScreen()
				return InitiateQueuesTab(&hndlr), cmd
			} else if tab.num == 2 {
				ClearScreen()
				return InitiateAddDownloadTab(&hndlr), cmd
			}
		case tea.KeyShiftRight:
			if tab.num == 1 {
				ClearScreen()
				return InitiateDownloadsTab(&hndlr), cmd
			} else if tab.num == 2 {
				ClearScreen()
				return InitiateQueuesTab(&hndlr), cmd
			} else if tab.num == 3 {
				ClearScreen()
				return InitiateAddDownloadTab(&hndlr), cmd
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return tab, tea.Quit
		}
	}

	if tab.num == 1 {

		// add download tab listeners and action handlers
		var addDownloadTab = tab.TAB.(AddDownloadTab)
		switch msg := msg.(type) {

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
		if queuesTab.action == "list" {
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
				case tea.KeyCtrlD:
					queuesTab.action = "new"
					tab.TAB = queuesTab
					return tab, nil
				case tea.KeyCtrlE:
					queuesTab.action = "edit"
				case tea.KeyCtrlA:
					queuesTab.action = "delete"
				}
			}
		} else if queuesTab.action == "delete" {

		} else if queuesTab.action == "edit" {

		} else if queuesTab.action == "new" {
			if queuesTab.name == "" {
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						queuesTab.name = queuesTab.nameInput.Value()
					}
					queuesTab.nameInput, cmd = queuesTab.nameInput.Update(msg)
				}
			} else if queuesTab.savePath == "" {
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						queuesTab.savePath = queuesTab.savePathInput.Value()
					}
					queuesTab.savePathInput, cmd = queuesTab.savePathInput.Update(msg)
				}
			} else if queuesTab.maximumDownloads == 0 {
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						maximumDownloads, err := strconv.Atoi(queuesTab.maximumDownloadsInput.Value())
						if err != nil {
							queuesTab.maximumDownloadsInput.Reset()
						}
						queuesTab.maximumDownloads = maximumDownloads
					}
					queuesTab.maximumDownloadsInput, cmd = queuesTab.maximumDownloadsInput.Update(msg)
				}
			} else if queuesTab.maximumBandWidth == 0 {
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						maximumBandWidth, err := strconv.ParseFloat(queuesTab.maximumBandWidthInput.Value(), 64)
						if err != nil {
							queuesTab.maximumBandWidthInput.Reset()
						}
						queuesTab.maximumBandWidth = maximumBandWidth
					}
					queuesTab.maximumBandWidthInput, cmd = queuesTab.maximumBandWidthInput.Update(msg)
				}
			} else if queuesTab.startTime == "" {
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.startTimeInput.Value() != "" {
							_, err := time.Parse("2006-02-01 15:04", queuesTab.startTimeInput.Value()[:16])
							if err != nil {
								queuesTab.startTimeInput.Reset()
							}
							queuesTab.startTime = queuesTab.startTimeInput.Value()
						} else {
							queuesTab.startTime = " "
						}

					}
					queuesTab.startTimeInput, cmd = queuesTab.startTimeInput.Update(msg)
				}
			} else if queuesTab.endTime == "" {
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.endTimeInput.Value() != "" {
							_, err := time.Parse("2006-02-01 15:04", queuesTab.endTimeInput.Value()[:16])
							if err != nil {
								queuesTab.endTimeInput.Reset()
							}
							queuesTab.endTime = queuesTab.endTimeInput.Value()
						} else {
							queuesTab.endTime = " "
						}
						var date entity.TimeInterval
						if queuesTab.startTime == " " {
							date.StartTime, _ = time.Parse("2006-02-01 15:04", "0001-01-01 00:00")
						} else {
							startTime, err := time.Parse("2006-02-01 15:04", queuesTab.startTime)
							if err != nil {
								panic(err)
							}
							date.StartTime = startTime
						}
						if queuesTab.endTime == " " {
							date.EndTime, _ = time.Parse("2006-02-01 15:04", "9999-30-12 23:59")
						} else {
							endTime, err := time.Parse("2006-02-01 15:04", queuesTab.endTime)
							if err != nil {
								panic(err)
							}
							date.EndTime = endTime
						}
						err := CreateQueue(queuesTab.name, queuesTab.savePath, queuesTab.maximumDownloads, queuesTab.maximumBandWidth, date)
						if err != nil {
							queuesTab.err = err
						}
						queuesTab.action = "finishedCreating"
					}
					queuesTab.endTimeInput, cmd = queuesTab.endTimeInput.Update(msg)
				}
			}
		} else if queuesTab.action[:8] == "finished" {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if key := msg.String(); strings.ToLower(key)[0] == 'y' {
					return InitiateDownloadsTab(&hndlr), nil
				} else if strings.ToLower(key)[0] == 'n' {
					return tab, tea.Quit
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
		view = "                   ----------------------------------------------------------- Add Download Tab -----------------------------------------------------------"
		var addDownloadTab = tab.TAB.(AddDownloadTab)
		if addDownloadTab.finished {
			view += "\nDownload added successfully!\n\nDo you want to continue? (y => yes) (n => no) (ctrl+c => quit)"
		} else if addDownloadTab.url == "" {
			view += fmt.Sprintf(
				"\nEnter the url here:\n\n%s\n\n%s",
				addDownloadTab.urlInput.View(),
				"(ctrl+c => quit)",
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
			view += "\n (ctrl+c => quit)"
		} else if addDownloadTab.fileName == "" {
			view += fmt.Sprintf(
				"\nEnter the file name here (optional):\n\n%s\n\n%s",
				addDownloadTab.fileNameInput.View(),
				"(ctrl+c => quit)",
			) + "\n"
		}
	} else if tab.num == 2 {
		var downloadsTab = tab.TAB.(DownloadsTab)
		view = "                   ----------------------------------------------------------- Downloads Tab -----------------------------------------------------------"
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
				view = fmt.Sprintf("%v\n%vURL: %v    Queue: %v    Status: %v    Speed: %vKB/s    Progress: %v", view, cursor, download.URL, queue.Name, download.Status, download.CurrentSpeed, download.Progress)
			}
			view = fmt.Sprintf("%v\n\n(ctrl+a => pause/resume) (ctrl+s => retry) (ctrl+d => delete) (ctrl+c => quit)", view)
		}
		if downloadsTab.deleteAction {
			view = fmt.Sprintf("%v\n%v\n%v", view, downloadsTab.message, "do you want to continue? (y => yes) (n => no) (ctrl+c => quit)")
		}
	} else if tab.num == 3 {
		var queuesTab = tab.TAB.(QueuesTab)
		view = "                   ----------------------------------------------------------- Queues Tab -----------------------------------------------------------"
		if queuesTab.action == "list" {
			view = fmt.Sprintf("%v\nSelect a queue:", view)
			for i, queue := range queuesTab.queues {
				cursor := "  "
				if i == queuesTab.cursorIndex {
					cursor = "> "
				}
				view = fmt.Sprintf("%v\n%vName: %v    Save-path: %v    Maximum-concurrent-downloads: %v    Maximum-band-width: %vKB/s    Activity-interval: from %v to %v", view, cursor, queue.Name, queue.SavePath, queue.MaximumDownloads, queue.MaximumBandWidth, queue.ActivityInterval.StartTime.Format("2006-01-02 15:04"), queue.ActivityInterval.EndTime.Format("2006-01-02 15:04"))
			}
			view = fmt.Sprintf("%v\n\n(ctrl+d => create) (ctrl+e => edit) (ctrl+a => delete) (ctrl+c => quit)", view)
		} else if queuesTab.action == "new" {
			if queuesTab.name == "" {
				view += fmt.Sprintf(
					"\nEnter the name of queue here:\n\n%s\n\n%s",
					queuesTab.nameInput.View(),
					"(ctrl+c => quit)",
				) + "\n"
			} else if queuesTab.savePath == "" {
				view += fmt.Sprintf(
					"\nEnter the path of queue here:\n\n%s\n\n%s",
					queuesTab.savePathInput.View(),
					"(ctrl+c => quit)",
				) + "\n"
			} else if queuesTab.maximumDownloads == 0 {
				view += fmt.Sprintf(
					"\nEnter the number of maximum downloads here:\n\n%s\n\n%s",
					queuesTab.maximumDownloadsInput.View(),
					"(ctrl+c => quit)",
				) + "\n"
			} else if queuesTab.maximumBandWidth == 0 {
				view += fmt.Sprintf(
					"\nEnter the maximum band width here:\n\n%s\n\n%s",
					queuesTab.maximumBandWidthInput.View(),
					"(ctrl+c => quit)",
				) + "\n"
			} else if queuesTab.startTime == "" {
				view += fmt.Sprintf(
					"\nEnter the start time of the queue here with format \"yyyy-dd-MM hh:mm\" (optional):\n\n%s\n\n%s",
					queuesTab.startTimeInput.View(),
					"(ctrl+c => quit)",
				) + "\n"
			} else if queuesTab.endTime == "" {
				view += fmt.Sprintf(
					"\nEnter the end time of the queue here with format \"yy-dd-MM hh:mm\" (optional):\n\n%s\n\n%s",
					queuesTab.endTimeInput.View(),
					"(ctrl+c => quit)",
				) + "\n"
			}
		} else if queuesTab.action[:8] == "finished" {
			var addition string
			var condition = queuesTab.action[8:len(queuesTab.action)]
			if condition == "Creating" {
				addition = "Queue added successfully"
			} else if condition == "Editing" {
				addition = "Queue edited successfully"
			} else if condition == "Deleting" {
				addition = "Queue deleted successfully"
			}
			view += fmt.Sprintf("\n%v\n\nDo you want to continue? (y => yes) (n => no) (ctrl+c => quit)", addition)
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

func CreateQueue(name string, savePath string, maximumDownloads int, maximumBandWidth float64, activityInterval entity.TimeInterval) error {
	return hndlr.QueueHandler.CreateQueue(dto.QueueDto{
		Name:             name,
		SavePath:         savePath,
		MaximumDownloads: maximumDownloads,
		MaximumBandWidth: maximumBandWidth,
		ActivityInterval: activityInterval,
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
