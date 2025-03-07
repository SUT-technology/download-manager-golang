package db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SUT-technology/download-manager-golang/internal/interface/config"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

type Pool struct {
	downloadPath string
	queuePath    string
}

func New(cfg config.DB) (*Pool, error) {
	return &Pool{
		downloadPath: cfg.Downloads,
		queuePath:    cfg.Queues,
	}, nil
}

// Query allows for querying data from multiple repositories (e.g., downloads and products).
func (p *Pool) Query(f func(r *repository.Repo) error) error {
	repo := &repository.Repo{
		Tables: repository.Tables{
			Downloads: newdownloadsTable(p),
			Queues:    newqueuesTable(p),
		},
	}

	// Execute the provided query function.
	return f(repo)
}

// Close simulates closing any resources (if needed).
func (p *Pool) Close() error {
	// In case of actual DB connections, you would close them here.
	return nil
}

func (p *Pool) loadData(filePath string, dst interface{}) error {
	files, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open downloads/queues file: %v", err)
	}
	defer files.Close()
	if err := json.NewDecoder(files).Decode(&dst); err != nil {
		return fmt.Errorf("could not decode data from file: %v", err)
	}

	return nil
}

func (p *Pool) saveData(filePath string, src interface{}) error {
	files, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create downloads/queues file: %v", err)
	}
	defer files.Close()
	if err := json.NewEncoder(files).Encode(src); err != nil {
		return fmt.Errorf("could not encode data to file: %v", err)
	}

	return nil
}
