package gormutil

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

type textFormatter struct{}

func (f *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	level := strings.ToUpper(entry.Level.String()[:4])

	fmt.Fprintf(b, "[%s][%s] ", level, entry.Time.Format("2006-01-02T15:04:05.000"))

	b.WriteString(entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func init() {
	logrus.SetFormatter(&textFormatter{})
}

func TestLogger(t *testing.T) {
	type Example struct {
		ID   int64
		Name string `gorm:"column:name;type:varchar(255);not null;default:''"`
	}
	c := MemoryOption()
	c.Logger = LoggerConfig{Logger: LoggerLogrus, LogLevel: "info"}
	db, err := Open(c)
	if err != nil {
		t.Error(err)
	}
	if err := db.AutoMigrate(&Example{}); err != nil {
		t.Error(err)
	}

	e := &Example{Name: "Dany"}
	if err := db.Create(e).Error; err != nil {
		t.Error(err)
	}
	t.Log(e.ID)
}
