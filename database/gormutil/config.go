package gormutil

import (
	"strconv"
	"time"

	"github.com/go-chocolate/configuration/common"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type Duration string

func (d Duration) value() time.Duration {
	if val, err := strconv.Atoi(string(d)); err == nil {
		return time.Duration(val) * time.Millisecond
	}
	dur, _ := time.ParseDuration(string(d))
	return dur
}

type Option = common.Config

type Config struct {
	Driver          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime Duration
	ConnMaxLifetime Duration
	Logger          LoggerConfig
	Option          Option
}

type LogLevel string

func (l LogLevel) Level() logger.LogLevel {
	switch l {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	}
	return logger.Info
}

const (
	LoggerStd    = "std"
	LoggerLogrus = "logrus"
)

type LoggerConfig struct {
	Logger                    string   // std logrus
	SlowThreshold             string   //慢查询定义，格式：2s 1s 200ms
	Colorful                  bool     //
	IgnoreRecordNotFoundError bool     //忽略 NotFoundError
	ParameterizedQueries      bool     //隐藏查询参数
	LogLevel                  LogLevel //日志打印级别 1 Silent, 2 Error, 3 Warn, 4 Info
}

func (l LoggerConfig) build() logger.Interface {
	config := logger.Config{
		Colorful:                  l.Colorful,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
		ParameterizedQueries:      l.ParameterizedQueries,
		LogLevel:                  l.LogLevel.Level(),
	}
	if v, err := time.ParseDuration(l.SlowThreshold); err == nil {
		config.SlowThreshold = v
	}
	switch l.Logger {
	case LoggerLogrus:
		return LogrusLogger(config)
	default:
		return logger.New(logrus.StandardLogger(), config)
	}
}

func MemoryOption() Config {
	return Config{
		Driver: SQLITE,
		Option: Option{"Database": ":memory:"},
	}
}
