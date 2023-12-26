package kv

import "github.com/go-chocolate/configuration/common"

type Option = common.Config

type Config struct {
	Driver string
	Option Option
}
