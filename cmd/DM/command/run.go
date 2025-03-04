package command

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/SUT-technology/download-manager-golang/internal/application/services"
	"github.com/SUT-technology/download-manager-golang/internal/infrastructure/db"
	"github.com/SUT-technology/download-manager-golang/internal/interface/config"
	"github.com/SUT-technology/download-manager-golang/internal/interface/handlers"
	"github.com/SUT-technology/download-manager-golang/internal/ui"
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
	slog.Debug("initialized json database")

	srvcs := services.New(db)

	// RUN UI AND USE IT

	var wg sync.WaitGroup

	go ui.Run(&wg)


	// ui.Run()

	// SAMPLE USE HANDLERS
	hndlrs := handlers.New(srvcs)

	// downloads, err := hndlrs.DownloadHndlr.GetDownloads()

	// if err != nil {
	// 	return err
	// }

	// fmt.Println(downloads)

	download , err := hndlrs.DownloadHndlr.GetDownloadById("3")
	if err != nil{
		return err
	}

	fmt.Println(download)



	wg.Wait()
	return nil
}
