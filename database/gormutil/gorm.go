package gormutil

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(c Config, options ...GORMOption) (*gorm.DB, error) {
	if c.Driver == "" {
		c.Driver = MYSQL
	}
	driver, ok := drivers[c.Driver]
	if !ok {
		return nil, fmt.Errorf("unknown database driver: %s", c.Driver)
	}
	dialer, err := driver(c.Option)
	if err != nil {
		return nil, fmt.Errorf("database driver on error :%s %v", c.Driver, err)
	}
	db, err := gorm.Open(dialer, applyOptions(&gorm.Config{Logger: logger.New(log.New(os.Stdout, "", log.LstdFlags), c.Logger.Build())}, options...))
	if err != nil {
		return nil, err
	}

	innerDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	innerDB.SetMaxOpenConns(c.MaxOpenConns)
	innerDB.SetMaxIdleConns(c.MaxIdleConns)
	innerDB.SetConnMaxIdleTime(c.ConnMaxIdleTime.value())
	innerDB.SetConnMaxLifetime(c.ConnMaxLifetime.value())
	return db, nil
}

func MustOpen(c Config, options ...GORMOption) *gorm.DB {
	db, err := Open(c, options...)
	if err != nil {
		panic(err)
	}
	return db
}
