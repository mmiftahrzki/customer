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
	if app.Run(); err != nil {
		logger.Panic(err)
	}

	// router.HandleFunc("GET /api/top-secert", nil)

	// auth := auth_pkg.New()
	// customerValidation := validation.New()

	// signUp := router_pkg.Endpoint{Path: "/api/auth/signup", Method: http.MethodPost}
	// signIn := router_pkg.Endpoint{Path: "/api/auth/signin", Method: http.MethodPost}

	// getToken := router_pkg.Endpoint{Path: "/api/auth/token", Method: http.MethodPost}

	// router.Handle(signUp, handler.CreateUser)
	// router.Handle(signIn, handler.ReadUser)
	// router.Handle(getToken, auth_pkg.Token)
}
