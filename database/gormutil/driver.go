package gormutil

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MYSQL  = "mysql"
	SQLITE = "sqlite"
)

type Driver func(c Option) (gorm.Dialector, error)

var drivers = map[string]Driver{
	MYSQL:  mysqlDriver,
	SQLITE: sqliteDriver,
}

func mysqlDriver(c Option) (gorm.Dialector, error) {
	u := url.URL{}
	//u.Scheme = MYSQL
	if username := c.String("Username"); username != "" {
		u.User = url.UserPassword(username, c.String("Password"))
	}
	v := url.Values{}
	u.Host = fmt.Sprintf("tcp(%s)", c.String("Addr"))
	u.Path = c.String("Database")
	v.Set("parseTime", "true")
	v.Set("loc", "Asia/Shanghai")
	u.RawQuery = v.Encode()
	dsn := strings.TrimPrefix(u.String(), "//")
	return mysql.Open(dsn), nil
}

func sqliteDriver(c Option) (gorm.Dialector, error) {
	dsn := c.String("Database")
	return sqlite.Open(dsn), nil
}

func RegisterDriver(name string, driver Driver) {
	drivers[name] = driver
}
