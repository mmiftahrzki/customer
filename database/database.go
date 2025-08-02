package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
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

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err.Error())

		return nil, fmt.Errorf("failed to open a connection to %s", cfg.Host)
	}

	db.SetMaxOpenConns(cfg.MaxConnection)
	db.SetMaxIdleConns(cfg.MaxConnection / 2)

	if err = db.Ping(); err != nil {
		log.Fatalln(err.Error())

		return nil, fmt.Errorf("failed to connect to %s", cfg.Host)
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
