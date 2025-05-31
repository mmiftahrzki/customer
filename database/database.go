package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/config"
)

var once sync.Once
var db_conn *sql.DB

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
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)

	return db, nil
}

func GetDatabaseConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	var err error

	if db_conn != nil {
		return db_conn, nil
	}

	once.Do(func() {
		db_conn, err = new(cfg)
	})

	return db_conn, err
}
