package dto

import "github.com/SUT-technology/download-manager-golang/internal/domain/entity"

type QueueDto struct {
	Name             string
	SavePath         string
	MaximumDownloads int
	MaximumBandWidth float64
	ActivityInterval entity.TimeInterval
}
