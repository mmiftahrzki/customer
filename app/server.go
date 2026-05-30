package app

import (
	"context"
	"database/sql"
	"errors"
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
			Handler:      newMux(db, cfg),
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 10,
		},
	}
}

func (a *app) Run() {
	a.log.Infof("listening on %s", a.server.Addr)

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Fatalln(err)
	} else {
		a.log.Infoln("application stopped gracefully")
	}
}

func (a *app) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Fatalf("shutdown failed: %v", err)
	} else {
		a.log.Infoln("application shutdown")
	}
}
