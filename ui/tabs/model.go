package tabs

import (
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/charmbracelet/bubbles/textinput"
)

type (
	errMsg error
)

var hndlr Handlers

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
	progress        int
	progressChan    chan int
}

type DownloadsTab struct {
	downloads    []entity.Download
	cursorIndex  int
	deleteAction bool
	message      string
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
	id                    string
	name                  string
	savePath              string
	maximumDownloads      int
	maximumBandWidth      float64
	startTime             string
	endTime               string
	message               string
	err                   error
}
