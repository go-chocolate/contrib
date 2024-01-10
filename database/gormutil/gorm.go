package gormutil

import (
	"fmt"

	"gorm.io/gorm"
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
	db, err := gorm.Open(dialer, applyOptions(&gorm.Config{Logger: c.Logger.build()}, options...))
	if err != nil {
		return nil, err
	}

	innerDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if c.MaxOpenConns > 0 {
		innerDB.SetMaxOpenConns(c.MaxOpenConns)
	}
	if c.MaxIdleConns > 0 {
		innerDB.SetMaxIdleConns(c.MaxIdleConns)
	}
	if v := c.ConnMaxIdleTime.value(); v > 0 {
		innerDB.SetConnMaxIdleTime(v)
	}
	if v := c.ConnMaxLifetime.value(); v > 0 {
		innerDB.SetConnMaxLifetime(v)
	}
	return db, nil
}

func MustOpen(c Config, options ...GORMOption) *gorm.DB {
	db, err := Open(c, options...)
	if err != nil {
		panic(err)
	}
	return db
}

func OpenMemory(options ...GORMOption) (*gorm.DB, error) {
	return Open(MemoryOption(), options...)
}
