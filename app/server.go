package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/mmiftahrzki/customer/config"
)

type app struct {
	server *http.Server
}

func New(cfg config.AppConfig, db *sql.DB) *app {
	return &app{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      newMux(db),
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 10,
		},
	}
}

func (a *app) Run() error {
	fmt.Println("Listening on:", a.server.Addr)

	return a.server.ListenAndServe()
}
