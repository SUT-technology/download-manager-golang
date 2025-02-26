package command

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/SUT-technology/download-manager-golang/internal/application/services"
	"github.com/SUT-technology/download-manager-golang/internal/infrastructure/db"
	"github.com/SUT-technology/download-manager-golang/internal/interface/config"
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers"
	"github.com/SUT-technology/download-manager-golang/pkg/tools/slogger"
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
	slog.Debug("initialized oracle database")

	srvcs := services.New(db)

	// RUN UI AND USE IT

	// SAMPLE USE HANDLERS
	hndlrs := handlers.New(srvcs)

	downloads, err := hndlrs.DownloadHndlr.GetDownloads()

	if err != nil {
		return err
	}

	fmt.Println(downloads)

	return nil
}
