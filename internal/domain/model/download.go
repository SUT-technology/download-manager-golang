package model

type DownloadProgressMsg struct {
	DownloadID string
	Progress   float64
	Speed      float64
	Status     string
	Downloaded int64
}

type DownloadControlMessage string

const (
	PauseCommand  DownloadControlMessage = "pause"
	ResumeCommand DownloadControlMessage = "resume"
)
