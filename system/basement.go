package system

import (
	"github.com/go-chocolate/contrib/database/gormutil"
	"github.com/go-chocolate/contrib/kv"
	"gorm.io/gorm"
)

type Basement struct {
	Database  *gorm.DB
	KVStorage kv.Storage
}

func (b *Basement) Setup(c Config) error {
	if db, err := gormutil.Open(c.Database); err != nil {
		return err
	} else {
		b.Database = db
	}

	if storage, err := kv.New(c.KVStorage); err != nil {
		return err
	} else {
		b.KVStorage = storage
	}

	return nil
}
