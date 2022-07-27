package storage

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"gorm.io/gorm"
)

type Interface interface {
	DB() *gorm.DB
	GetUserCenter(account string) (*models.UserCenter, error)
	Good(barcodes string, account string) (*models.TemporaryGoodsInventories, error)
}
