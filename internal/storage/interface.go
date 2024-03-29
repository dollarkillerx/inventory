package storage

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"gorm.io/gorm"
)

type Interface interface {
	DB() *gorm.DB

	GetUserCenter(account string) (*models.UserCenter, error)
	Good(barcodes string, account string) (*models.TemporaryGoodsInventories, error)
	DeleteGood(goodID string, account string) (err error)
	Search(keyword string, account string) ([]models.TemporaryGoodsInventories, error)
	WareHousing(goodsId string, barcode string, account string, cost float64, numberProducts int, remark string) (err error)
	OutStock(goodsId string, barcode string, account string, cost float64, numberProducts int, price float64, remark string) (err error)
	IOHistory(goodsID string, account string) ([]models.InventoryHistoryDetailed, error)
	IORevoke(orderID string, account string) ([]models.InventoryHistoryDetailed, error)
	IOList(account string, limit int, offset int) (count int64, ids []models.InventoryHistoryDetailed, err error)
	ResetStatistics() (err error)

	// Statistics ..
	Statistics(account string) (scs []models.StatisticsBand, err error)
}
