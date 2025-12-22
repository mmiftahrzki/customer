package database

import (
	"database/sql"
	"fmt"
	"sync"

	asd "github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/config"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var db *sql.DB
var log *logrus.Entry = logger.GetLogger().WithField("component", "database")

func new(cfg config.DatabaseConfig) (*sql.DB, error) {
	var err error
	var db *sql.DB
	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	err = asd.SetLogger(log)
	if err != nil {
		return nil, fmt.Errorf("failed to set Logger: %w", err)
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open a connection to %s: %w", cfg.Host, err)
	}

	db.SetMaxOpenConns(cfg.MaxConnection)
	db.SetMaxIdleConns(cfg.MaxConnection / 2)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", cfg.Host, err)
	}

	return db, nil
}

func New(cfg config.DatabaseConfig) (*sql.DB, error) {
	var err error

	if db != nil {
		return db, nil
	}

	once.Do(func() {
		db, err = new(cfg)
	})

	return db, err
}
