package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/app"
	"github.com/mmiftahrzki/customer/config"
	"github.com/mmiftahrzki/customer/database"
	"github.com/mmiftahrzki/customer/logger"
)

func main() {
	logger := logger.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalln(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		logger.Fatalf("Database Error: %v\n", err)
	}
	defer db.Close()

	app := app.New(cfg.App, db)

	go app.Run()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopCh
	defer close(stopCh)

	logger.Infoln("Shutting down server...")

	app.Shutdown(context.Background())
}
