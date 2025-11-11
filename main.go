package main

import (
	_ "embed"

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
	err = app.Run()
	if err != nil {
		logger.Panic(err)
	}
}
