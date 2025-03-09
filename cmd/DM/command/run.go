package command

import (
	"flag"
	"fmt"
	"github.com/SUT-technology/download-manager-golang/internal/application/services"
	"github.com/SUT-technology/download-manager-golang/internal/domain/dto"
	"github.com/SUT-technology/download-manager-golang/internal/domain/entity"
	"github.com/SUT-technology/download-manager-golang/internal/infrastructure/db"
	"github.com/SUT-technology/download-manager-golang/internal/interface/config"
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/slogger"
	"log/slog"
	"os"
)

func Run() error {
	var configPath string
	flag.StringVar(&configPath, "cfg", "assets/config/development.yaml", "Configuration File")
	flag.Parse()
	c, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	logger := slogger.NewJSONLogger(c.Logger.Level, os.Stdout)
	slog.SetDefault(logger)

	db, err := db.New(c.DB)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	slog.Debug("initialized json database")

	srvcs := services.New(db)

	// RUN UI AND USE IT

	// SAMPLE USE HANDLERS
	hndlrs := handlers.New(srvcs)

	// downloads, err := hndlrs.DownloadHndlr.GetDownloads()

	// if err != nil {
	// 	return err
	// }

	// fmt.Println(downloads)

	hndlrs.QueueHndlr.CreateQueue(dto.QueueDto{
		"tst22",
		"tmp/tmp2",
		0,
		0,
		entity.TimeInterval{},
	})
	queue, err := hndlrs.QueueHndlr.GetQueueById("2")
	if err != nil {
		return fmt.Errorf("getting queue: %w", err)
	}

	hndlrs.DownloadHndlr.CreateDownload(dto.DownloadDto{
		URL:      "https://dl.nakaman-music.ir/Music/BAHRAM/Forsat/Bahram%20-%20Gear%20Box.mp3",
		Queue:    queue,
		FileName: "bahram.mp3",
	})

	return nil
}
