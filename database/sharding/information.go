package sharding

import (
	"bytes"
	"fmt"
	"time"
)

type Information struct {
	ID          int64     `gorm:""`
	Name        string    `gorm:"type:varchar(255);not null;index"`
	Table       string    `gorm:"type:varchar(255);not null;index"`
	CreatedTime time.Time `gorm:"type:datetime"`
	UpdatedTime time.Time `gorm:"type:datetime"`
}

func (i Information) tableName() string {
	return fmt.Sprintf("%s_sharding", i.Name)
}

func (i *Information) exec() []string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("CREATE TABLE IF NOT EXISTS")
	buf.WriteString("`" + i.tableName() + "`")
	buf.WriteString("(")
	buf.WriteString("`ID` BIGINT NOT NULL,")
	buf.WriteString("`NAME` VARCHAR(255) NOT NULL,")
	buf.WriteString("`TABLE` VARCHAR(255) NOT NULL,")
	buf.WriteString("`CREATED_TIME` datetime,")
	buf.WriteString("`UPDATED_TIME` datetime,")
	buf.WriteString("PRIMARY KEY (`ID`),")
	buf.WriteString("INDEX " + "`IDX_" + i.tableName() + "_NAME`(`NAME`),")
	buf.WriteString("UNIQUE INDEX `IDX_" + i.tableName() + "_TABLE`(`TABLE`)")
	buf.WriteString(")")
	// WARNING: sqlite不支持这种INDEX写法
	exec := []string{
		buf.String(),
	}
	return exec
}
