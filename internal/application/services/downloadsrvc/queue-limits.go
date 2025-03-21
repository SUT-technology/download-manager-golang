package downloadsrvc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/repository"
)

var activeDownloadsMap = make(map[string][]string)
var activeDownloadsMutex sync.Mutex

func (d DownloadService) GetQueueById(id string) (*entity.Queue, error) {
	var (
		queue *entity.Queue
		err   error
	)

	ctx := context.Background()

	queryFunc := func(r *repository.Repo) error {
		queue, err = r.Tables.Queues.GetQueueById(ctx, id)
		if err != nil {
			return fmt.Errorf("getting data from queues: %w", err)
		}

		return nil
	}

	err = d.db.Query(queryFunc)
	if err != nil {
		return nil, err
	}

	return queue, nil
}
func incrementActiveDownload(queueID string, downloadID string) {
	activeDownloadsMutex.Lock()
	seen := make(map[string]bool)
	for _, v := range activeDownloadsMap[queueID] {
		seen[v] = true
	}

	if !seen[downloadID] {
		activeDownloadsMap[queueID] = append(activeDownloadsMap[queueID], downloadID)
	}
	activeDownloadsMutex.Unlock()
}

func decrementActiveDownload(queueID string, downloadID string) {
	activeDownloadsMutex.Lock()
	for i, v := range activeDownloadsMap[queueID] {
		if v == downloadID {
			activeDownloadsMap[queueID] = append(activeDownloadsMap[queueID][:i], activeDownloadsMap[queueID][i+1:]...)
			break
		}
	}
	activeDownloadsMutex.Unlock()
}

func GetActiveDownload(queueID string) int {
	activeDownloadsMutex.Lock()
	defer activeDownloadsMutex.Unlock()
	return len(activeDownloadsMap[queueID])
}

func waitForAllowedTime(q entity.Queue) {
	now := time.Now()
	year, mon, day := now.Date()
	loc := now.Location()

	// Build today's allowed start and end times based solely on time-of-day.
	allowedStart := time.Date(year, mon, day,
		q.ActivityInterval.StartTime.Hour(),
		q.ActivityInterval.StartTime.Minute(),
		q.ActivityInterval.StartTime.Second(), 0, loc)
	allowedEnd := time.Date(year, mon, day,
		q.ActivityInterval.EndTime.Hour(),
		q.ActivityInterval.EndTime.Minute(),
		q.ActivityInterval.EndTime.Second(), 0, loc)

	// If before allowed start, sleep until allowedStart.
	if now.Before(allowedStart) {
		time.Sleep(allowedStart.Sub(now))
		return
	}
	// If after allowed end, sleep until next day's allowed start.
	if now.After(allowedEnd) {
		nextAllowedStart := allowedStart.Add(24 * time.Hour)
		time.Sleep(nextAllowedStart.Sub(now))
		return
	}
	// Otherwise we are within the allowed window so return immediately.
}
