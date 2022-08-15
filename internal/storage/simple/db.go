package simple

import (
	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/dollarkillerx/inventory/internal/utils"
	"gorm.io/gorm"
	"sync"
)

type Simple struct {
	db *gorm.DB

	inventoryMu sync.Mutex
}

func NewSimple(conf *conf.PgSQLConfig) (*Simple, error) {
	sql, err := utils.InitPgSQL(conf)
	if err != nil {
		return nil, err
	}

	sql.AutoMigrate(
		&models.UserCenter{},
		//&models.Goods{},
		&models.History{},
		&models.Inventory{},
		&models.InventoryModify{},
		&models.InventoryHistory{},
		&models.InventoryHistoryDetailed{},
	)

	return &Simple{
		db: sql,
	}, nil
}

func (s *Simple) DB() *gorm.DB {
	return s.db
}
