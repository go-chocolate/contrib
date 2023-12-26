package gormutil

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type GORMOption func(o *gorm.Config)

func applyOptions(o *gorm.Config, options ...GORMOption) *gorm.Config {
	for _, opt := range options {
		opt(o)
	}
	return o
}

func WithSkipDefaultTransaction(SkipDefaultTransaction bool) GORMOption {
	return func(o *gorm.Config) {
		o.SkipDefaultTransaction = SkipDefaultTransaction
	}
}

func WithNamingStrategy(NamingStrategy schema.Namer) GORMOption {
	return func(o *gorm.Config) {
		o.NamingStrategy = NamingStrategy
	}
}

func WithFullSaveAssociations(FullSaveAssociations bool) GORMOption {
	return func(o *gorm.Config) {
		o.FullSaveAssociations = FullSaveAssociations
	}
}

func WithLogger(Logger logger.Interface) GORMOption {
	return func(o *gorm.Config) {
		o.Logger = Logger
	}
}

func WithNowFunc(NowFunc func() time.Time) GORMOption {
	return func(o *gorm.Config) {
		o.NowFunc = NowFunc
	}
}

func WithDryRun(DryRun bool) GORMOption {
	return func(o *gorm.Config) {
		o.DryRun = DryRun
	}
}

func WithPrepareStmt(PrepareStmt bool) GORMOption {
	return func(o *gorm.Config) {
		o.PrepareStmt = PrepareStmt
	}
}

func WithDisableAutomaticPing(DisableAutomaticPing bool) GORMOption {
	return func(o *gorm.Config) {
		o.DisableAutomaticPing = DisableAutomaticPing
	}
}

func WithDisableForeignKeyConstraintWhenMigrating(DisableForeignKeyConstraintWhenMigrating bool) GORMOption {
	return func(o *gorm.Config) {
		o.DisableForeignKeyConstraintWhenMigrating = DisableForeignKeyConstraintWhenMigrating
	}
}

func WithIgnoreRelationshipsWhenMigrating(IgnoreRelationshipsWhenMigrating bool) GORMOption {
	return func(o *gorm.Config) {
		o.IgnoreRelationshipsWhenMigrating = IgnoreRelationshipsWhenMigrating
	}
}

func WithDisableNestedTransaction(DisableNestedTransaction bool) GORMOption {
	return func(o *gorm.Config) {
		o.DisableNestedTransaction = DisableNestedTransaction
	}
}

func WithAllowGlobalUpdate(AllowGlobalUpdate bool) GORMOption {
	return func(o *gorm.Config) {
		o.AllowGlobalUpdate = AllowGlobalUpdate
	}
}

func WithQueryFields(QueryFields bool) GORMOption {
	return func(o *gorm.Config) {
		o.QueryFields = QueryFields
	}
}

func WithCreateBatchSize(CreateBatchSize int) GORMOption {
	return func(o *gorm.Config) {
		o.CreateBatchSize = CreateBatchSize
	}
}

func WithTranslateError(TranslateError bool) GORMOption {
	return func(o *gorm.Config) {
		o.TranslateError = TranslateError
	}
}

func WithClauseBuilders(ClauseBuilders map[string]clause.ClauseBuilder) GORMOption {
	return func(o *gorm.Config) {
		o.ClauseBuilders = ClauseBuilders
	}
}

func WithConnPool(ConnPool gorm.ConnPool) GORMOption {
	return func(o *gorm.Config) {
		o.ConnPool = ConnPool
	}
}

func WithPlugins(Plugins map[string]gorm.Plugin) GORMOption {
	return func(o *gorm.Config) {
		o.Plugins = Plugins
	}
}

func WithStdLogger(level ...logger.LogLevel) GORMOption {
	return WithLogger(StdLogger(level...))
}

func StdLogger(level ...logger.LogLevel) logger.Interface {
	lev := logger.Info
	if len(level) > 0 {
		lev = level[0]
	}
	return logger.New(log.New(os.Stdout, "[GORM]", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Millisecond * 200,
		Colorful:                  false,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		LogLevel:                  lev,
	})
}
