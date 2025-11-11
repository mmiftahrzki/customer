package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/mmiftahrzki/customer/config"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

type app struct {
	server *http.Server
	log    *logrus.Entry
}

func New(cfg config.AppConfig, db *sql.DB) *app {
	app_logger := logger.GetLogger().WithField("component", "app")

	return &app{
		log: app_logger,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      newMux(db),
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 10,
		},
	}
}

func (a *app) Run() error {
	a.log.Infof("Listening on %s", a.server.Addr)

	return a.server.ListenAndServe()
}
