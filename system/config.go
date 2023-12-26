package system

import (
	"github.com/go-chocolate/contrib/database/gormutil"
	"github.com/go-chocolate/contrib/kv"
)

type Config struct {
	Database  gormutil.Config
	KVStorage kv.Config
}
