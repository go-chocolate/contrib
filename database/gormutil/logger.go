package gormutil

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type logrusLogger struct {
	logger.Config
}

var _ logger.Interface = (*logrusLogger)(nil)

func (l *logrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	copied := new(logrusLogger)
	*copied = *l
	copied.LogLevel = level
	return copied
}

func (l *logrusLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	logrus.WithContext(ctx).Infof(format, args...)
}

func (l *logrusLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	logrus.WithContext(ctx).Warnf(format, args...)
}

func (l *logrusLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	logrus.WithContext(ctx).Errorf(format, args...)
}

func (l *logrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	var exec string
	var rowsAffected int64
	log := logrus.WithContext(ctx)
	defer func() {
		if l.LogLevel == logger.Silent {
			return
		}
		if l.LogLevel >= logger.Warn {
			log.Infof("[ROWS:%d][%v] %s", rowsAffected, time.Since(begin), exec)
		}
		if err != nil {
			log.Errorf("%v", err)
		}
	}()
	exec, rowsAffected = fc()
}
