package entity

type Download struct {
	ID           string  `json:"id"`
	URL          string  `json:"url"`
	QueueId      string  `json:"queue_id"`
	FileName     string  `json:"fileName"`
	Status       int     `json:"status"`
	CurrentSpeed float64 `json:"currentSpeed"`
	Progress     float64 `json:"progress"`
}
