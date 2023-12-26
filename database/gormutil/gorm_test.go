package gormutil

import "testing"

func TestOpen_MySQL(t *testing.T) {
	db, err := Open(Config{
		Driver: MYSQL,
		Option: Option{
			"Addr":     "127.0.0.1:3306",
			"Username": "root",
			"Password": "root",
			"Database": "example",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	innerDB, err := db.DB()
	if err != nil {
		t.Error(err)
		return
	}
	defer innerDB.Close()
	if err := innerDB.Ping(); err != nil {
		t.Error(err)
	}
}

func TestOpen_SQLite(t *testing.T) {
	db, err := Open(Config{
		Driver: SQLITE,
		Option: Option{
			"Database": ":memory:",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	innerDB, err := db.DB()
	if err != nil {
		t.Error(err)
		return
	}
	defer innerDB.Close()
	if err := innerDB.Ping(); err != nil {
		t.Error(err)
	}
}
