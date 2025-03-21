package entity

type Download struct {
	ID           string  `json:"id"`
	URL          string  `json:"url"`
	QueueId      string  `json:"queue_id"`
	FileName     string  `json:"fileName"`
	TotalSize    int64   `json:"totalSize"`
	Downloaded   int64   `json:"downloaded"`
	Progress     float64 `json:"progress"`
	CurrentSpeed float64 `json:"currentSpeed"`
	Status       string  `json:"status"`
	OutPath       string  `json:"outPath"`
}


