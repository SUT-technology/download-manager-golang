package entity

type Queue struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	SavePath         string       `json:"savePath"`
	MaximumDownloads int          `json:"maximumDownloads"`
	MaximumBandWidth float64      `json:"maximumBandWidth"`
	ActivityInterval TimeInterval `json:"activityInterval"`
}
