package simple

import (
	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/dollarkillerx/inventory/internal/utils"
	"gorm.io/gorm"
)

type Simple struct {
	db *gorm.DB
}

func NewSimple(conf *conf.PgSQLConfig) (*Simple, error) {
	sql, err := utils.InitPgSQL(conf)
	if err != nil {
		return nil, err
	}

	sql.AutoMigrate(
		&models.UserCenter{},
		&models.Goods{},
	)

	return &Simple{
		db: sql,
	}, nil
}

func (s *Simple) DB() *gorm.DB {
	return s.db
}
