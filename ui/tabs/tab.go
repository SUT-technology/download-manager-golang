package tabs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SUT-technology/download-manager-golang/internal/application/services/downloadsrvc"
	"github.com/SUT-technology/download-manager-golang/internal/domain/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (tab Tab) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		WatchProgressCmd(),
	)
}

func (tab Tab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch message := msg.(type) {

	// general key listeners: shifting between tabs and quit
	case tea.KeyMsg:

		switch message.Type {
		case tea.KeyLeft:
			if tab.num == 3 {
				ClearScreen()
				return InitiateDownloadsTab(&hndlr), WatchProgressCmd()
			} else if tab.num == 1 {
				ClearScreen()
				return InitiateQueuesTab(&hndlr), cmd
			} else if tab.num == 2 {
				ClearScreen()
				return InitiateAddDownloadTab(&hndlr), cmd
			}
		case tea.KeyRight:
			if tab.num == 1 {
				ClearScreen()
				return InitiateDownloadsTab(&hndlr), WatchProgressCmd()
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
		// add download tab

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
			// get the url form input
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					addDownloadTab.url = addDownloadTab.urlInput.Value()
				}
				addDownloadTab.urlInput, cmd = addDownloadTab.urlInput.Update(msg)
			}

		} else if addDownloadTab.selectedQueueId == "" {
			// select a queue based on their name
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
			// get the file name form input
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.Type {
				case tea.KeyEnter:
					addDownloadTab.fileName = addDownloadTab.fileNameInput.Value()
					_, err := CreateDownload(addDownloadTab.url, addDownloadTab.selectedQueueId, addDownloadTab.fileName)
					if err != nil {
						addDownloadTab.err = err
					} else {
						return InitiateDownloadsTab(&hndlr), WatchProgressCmd()
					}
				}
				addDownloadTab.fileNameInput, cmd = addDownloadTab.fileNameInput.Update(msg)
			}
		}

		tab.TAB = addDownloadTab
		return tab, cmd

	} else if tab.num == 2 {
		//downloads list tab

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
					//pause/resume
					curDl := downloadsTab.downloads[downloadsTab.cursorIndex]
					controlChan, exists := downloadsrvc.ControlChannels[curDl.ID]
					if exists {
						if curDl.Status == "downloading" {
							controlChan <- model.PauseCommand
						} else if curDl.Status == "paused" {
							controlChan <- model.ResumeCommand
						}
					}

				case tea.KeyCtrlS:
					// Retry the selected download if not complete.
					curDl := downloadsTab.downloads[downloadsTab.cursorIndex]
					if curDl.Status != "completed" {
						curDl.Downloaded = 0
						curDl.Progress = 0
						curDl.Status = "downloading"
						downloadsTab.downloads[downloadsTab.cursorIndex] = curDl
						// Create a new control channel and start a new worker.
						controlChan := make(chan model.DownloadControlMessage)
						downloadsrvc.ControlChannels[curDl.ID] = controlChan
						go hndlr.DownloadHandler.DownloadWorker(&downloadsTab.downloads[downloadsTab.cursorIndex], downloadsrvc.ProgressChan, controlChan)
					}
				case tea.KeyCtrlD:
					// press ctrl+d to delete the selected download
					download, err := hndlr.DownloadHandler.DeleteDownload(downloadsTab.downloads[downloadsTab.cursorIndex].ID)
					if err != nil {
						panic(err)
					}
					downloadsTab.message = fmt.Sprintf("%s deleted successfully", download.FileName)
					downloadsTab.deleteAction = true
				}
			case model.DownloadProgressMsg:
				// Update the matching download’s progress and status.
				for i, d := range downloadsTab.downloads {
					if d.ID == msg.DownloadID {
						downloadsTab.downloads[i].Progress = msg.Progress
						downloadsTab.downloads[i].CurrentSpeed = msg.Speed
						downloadsTab.downloads[i].Status = msg.Status
						downloadsTab.downloads[i].Downloaded = msg.Downloaded
						break
					}
				}
				tab.TAB = downloadsTab
				return tab, WatchProgressCmd()
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
		// queues list tab

		queuesTab := tab.TAB.(QueuesTab)

		if queuesTab.action == "list" {
			// the list panel key handlings
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
				case tea.KeyCtrlW:
					queuesTab.action = "edit"
				case tea.KeyCtrlQ:
					queue, err := hndlr.QueueHandler.DeleteQueue(queuesTab.queues[queuesTab.cursorIndex].ID)
					if err != nil {
						panic(err)
					}
					// TEMPORARY - shall be handled in the database layers
					downloads, err := hndlr.DownloadHandler.GetDownloads()
					for _, download := range downloads {
						if download.QueueId == queue.ID {
							hndlr.DownloadHandler.DeleteDownload(download.ID)
						}
					}
					queuesTab.action = "finishedDeleting"
				}
			}

		} else if queuesTab.action == "edit" {
			// edit panel
			// initiating the edit panel textinputs placeholders
			queue := queuesTab.queues[queuesTab.cursorIndex]
			queuesTab.id = queue.ID
			queuesTab.nameInput.Placeholder = queue.Name
			queuesTab.savePathInput.Placeholder = queue.SavePath
			queuesTab.maximumDownloadsInput.Placeholder = fmt.Sprintf("%v", queue.MaximumDownloads)
			queuesTab.maximumBandWidthInput.Placeholder = fmt.Sprintf("%v", queue.MaximumBandWidth)

			tempStart := queue.ActivityInterval.StartTime.Format("15:04:05")
			if tempStart != "00:00:00" {
				queuesTab.startTimeInput.Placeholder = tempStart
			} else {
				queuesTab.startTimeInput.Placeholder = " "
			}

			tempEnd := queue.ActivityInterval.EndTime.Format("15:04:05")
			if tempEnd != "23:59:59" {
				queuesTab.endTimeInput.Placeholder = tempEnd
			} else {
				queuesTab.endTimeInput.Placeholder = " "
			}

			if queuesTab.name == "" {
				// get the new name of the queue from input
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.nameInput.Value() != "" {
							queuesTab.name = queuesTab.nameInput.Value()
						} else {
							queuesTab.name = queuesTab.nameInput.Placeholder
						}
					}
					queuesTab.nameInput, cmd = queuesTab.nameInput.Update(msg)
				}

			} else if queuesTab.savePath == "" {
				// get the new save path of the queue from input
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.savePathInput.Value() != "" {
							queuesTab.savePath = queuesTab.savePathInput.Value()
						} else {
							queuesTab.savePath = queuesTab.savePathInput.Placeholder
						}
					}
					queuesTab.savePathInput, cmd = queuesTab.savePathInput.Update(msg)
				}

			} else if queuesTab.maximumDownloads == 0 {
				// get the new maximum number of concurrent downloads from input
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.maximumDownloadsInput.Value() != "" {
							maximumDownloads, err := strconv.Atoi(queuesTab.maximumDownloadsInput.Value())
							if err != nil {
								queuesTab.maximumDownloadsInput.Reset()
							}
							queuesTab.maximumDownloads = maximumDownloads
						} else {
							queuesTab.maximumDownloads, _ = strconv.Atoi(queuesTab.maximumDownloadsInput.Placeholder)
						}
					}
					queuesTab.maximumDownloadsInput, cmd = queuesTab.maximumDownloadsInput.Update(msg)
				}

			} else if queuesTab.maximumBandWidth == 0 {
				// get the new maximum bandwidth from input
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.maximumBandWidthInput.Value() != "" {
							maximumBandWidth, err := strconv.ParseFloat(queuesTab.maximumBandWidthInput.Value(), 64)
							if err != nil {
								queuesTab.maximumBandWidthInput.Reset()
							}
							queuesTab.maximumBandWidth = maximumBandWidth
						} else {
							queuesTab.maximumBandWidth, _ = strconv.ParseFloat(queuesTab.maximumBandWidthInput.Placeholder, 64)
						}
					}
					queuesTab.maximumBandWidthInput, cmd = queuesTab.maximumBandWidthInput.Update(msg)
				}

			} else if queuesTab.startTime == "" {
				// get the new start time from input
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.startTimeInput.Value() != "" {
							_, err := time.Parse("15:04:05", queuesTab.startTimeInput.Value())
							if err != nil {
								queuesTab.startTimeInput.Reset()
							}
							queuesTab.startTime = queuesTab.startTimeInput.Value()
						} else {
							queuesTab.startTime = queuesTab.startTimeInput.Placeholder
						}
					}
					queuesTab.startTimeInput, cmd = queuesTab.startTimeInput.Update(msg)
				}

			} else if queuesTab.endTime == "" {
				// get the new end time from input
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.Type {
					case tea.KeyEnter:
						if queuesTab.endTimeInput.Value() != "" {
							_, err := time.Parse("15:04:05", queuesTab.endTimeInput.Value())
							if err != nil {
								// queuesTab.endTimeInput.Reset()
							}
							queuesTab.endTime = queuesTab.endTimeInput.Value()
						} else {
							queuesTab.endTime = queuesTab.startTimeInput.Placeholder
						}
						var date = getInterval(queuesTab.startTime, queuesTab.endTime)
						err := UpdateQueue(queue.ID, queuesTab.name, queuesTab.savePath, queuesTab.maximumDownloads, queuesTab.maximumBandWidth, date)
						if err != nil {
							queuesTab.err = err
						}
						queuesTab.action = "finishedEditing"
					}
					queuesTab.endTimeInput, cmd = queuesTab.endTimeInput.Update(msg)
				}

			}

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
							_, err := time.Parse("15:04:05", queuesTab.startTimeInput.Value())
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
							_, err := time.Parse("15:04:05", queuesTab.endTimeInput.Value())
							if err != nil {
								queuesTab.endTimeInput.Reset()
							}
							queuesTab.endTime = queuesTab.endTimeInput.Value()
						} else {
							queuesTab.endTime = " "
						}
						date := getInterval(queuesTab.startTime, queuesTab.endTime)
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
		view = "----------------------------------------------------------- Add Download Tab -----------------------------------------------------------"
		var addDownloadTab = tab.TAB.(AddDownloadTab)

		if addDownloadTab.err != nil {
			view += fmt.Sprintf("\n %v", addDownloadTab.err)
		}
		if addDownloadTab.url == "" {
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
		} else {
			// Show download progress if available
			if addDownloadTab.progress > 0 {
				view += fmt.Sprintf("\nDownload Progress: %d%%", addDownloadTab.progress)
			}
		}
	} else if tab.num == 2 {
		var downloadsTab = tab.TAB.(DownloadsTab)
		view = "----------------------------------------------------------- Downloads Tab -----------------------------------------------------------"
		if !downloadsTab.deleteAction {
			view = fmt.Sprintf("%v\nSelect a download:", view)
			for i, download := range downloadsTab.downloads {
				queue, err := hndlr.QueueHandler.GetQueueById(download.QueueId)
				if err != nil {
					panic(err)
				}
				cursor := "  "
				if i == downloadsTab.cursorIndex {
					cursor = "> "
				}
				view = fmt.Sprintf("%v\n%vQueue: %v    Status: %v    Speed: %vKB/s    Progress: %v%%", view, cursor, queue.Name, download.Status, int(download.CurrentSpeed), int(download.Progress))

			}
			view = fmt.Sprintf("%v\n\n(ctrl+a => pause/resume) (ctrl+s => retry) (ctrl+d => delete) (ctrl+c => quit)", view)
		}
		if downloadsTab.deleteAction {
			view = fmt.Sprintf("%v\n%v\n%v", view, downloadsTab.message, "do you want to continue? (y => yes) (n => no)")
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
			view = fmt.Sprintf("%v\n\n(ctrl+d => create) (ctrl+w => edit) (ctrl+q => delete) (ctrl+c => quit)", view)
		} else if queuesTab.action == "new" || queuesTab.action == "edit" {
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
					"\nEnter the start time of the queue here with format \"hh:mm:ss\" (optional):\n\n%s\n\n%s",
					queuesTab.startTimeInput.View(),
					"(ctrl+c => quit)",
				) + "\n"
			} else if queuesTab.endTime == "" {
				view += fmt.Sprintf(
					"\nEnter the end time of the queue here with format \"hh:mm:ss\" (optional):\n\n%s\n\n%s",
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
			view += fmt.Sprintf("\n%v\n\nDo you want to continue? (y => yes) (n => no)", addition)
		}
	}
	return view
}
