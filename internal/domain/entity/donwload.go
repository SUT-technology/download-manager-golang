package entity

import "net/url"

type Download struct {
	ID           string  `json:"id"`
	URL          url.URL `json:"url"`
	Queue        *Queue  `json:"queue"`
	FileName     string  `json:"fileName"`
	Status       int     `json:"status"`
	CurrentSpeed float64 `json:"currentSpeed"`
	Progress     float64 `json:"progress"`
}
