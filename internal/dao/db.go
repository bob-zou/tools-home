package dao

import (
	"database/sql"
	"time"
	"tools-home/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() (db *gorm.DB, cf func(), err error) {
	var cfg struct {
		Dsn             string
		ConnMaxLifetime time.Duration
		MaxOpenConn     int
		MaxIdleConn     int
	}

	if err = conf.Load("db.json", &cfg); err != nil {
		return
	}

	sqlDB, err := sql.Open("mysql", cfg.Dsn)
	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return
	}
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Second)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	if err = sqlDB.Ping(); err != nil {
		return
	}

	cf = func() { _ = sqlDB.Close() }
	return
}
