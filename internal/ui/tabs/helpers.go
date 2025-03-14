package tabs

import (
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"os"
	"os/exec"
	"runtime"
)

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
