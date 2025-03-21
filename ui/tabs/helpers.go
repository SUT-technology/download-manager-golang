package tabs

import (
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/SUT-technology/download-manager-golang/internal/application/services/downloadsrvc"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	tea "github.com/charmbracelet/bubbletea"
)

func CreateDownload(url, queueId, fileName string) (*entity.Download, error) {
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
func UpdateQueue(id string, name string, savePath string, maximumDownloads int, maximumBandWidth float64, activityInterval entity.TimeInterval) error {
	return hndlr.QueueHandler.FindAndUpdateQueue(id, dto.QueueDto{
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

func getInterval(startTime string, endTime string) (date entity.TimeInterval) {
	if startTime == " " {
		date.StartTime, _ = time.Parse("2006-02-01 15:04", "0001-01-01 00:00")
	} else {
		startTime, err := time.Parse("2006-02-01 15:04", startTime)
		if err != nil {
			panic(err)
		}
		date.StartTime = startTime
	}
	if endTime == " " {
		date.EndTime, _ = time.Parse("2006-02-01 15:04", "9999-30-12 23:59")
	} else {
		endTime, err := time.Parse("2006-02-01 15:04", endTime)
		if err != nil {
			panic(err)
		}
		date.EndTime = endTime
	}
	return date
}

func WatchProgressCmd() tea.Cmd {
	return func() tea.Msg {
		return <-downloadsrvc.ProgressChan
	}
}
