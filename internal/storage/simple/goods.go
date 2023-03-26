package simple

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"

	"log"
	"strings"
)

func (s *Simple) Goods() {

}

func (s *Simple) DeleteGood(goodID string, account string) (err error) {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 涉及到的关联表

		// 1. goods 商品表
		var good models.Goods
		err = tx.Model(&models.Goods{}).
			Where("id = ?", goodID).Where("by_account = ?", account).First(&good).Error
		if err != nil {
			return err
		}

		if good.ByAccount != account || good.ID != goodID {
			return errors.New("非法请求")
		}

		err = tx.Model(&models.Goods{}).
			Where("id = ?", goodID).
			//Where("by_account = ?", account).
			Unscoped().Delete(&models.Goods{}).Error
		if err != nil {
			return err
		}

		// 2. inventory 库存表
		err = tx.Model(&models.Inventory{}).
			Where("goods_id = ?", goodID).Unscoped().Delete(&models.Inventory{}).Error
		if err != nil {
			return err
		}

		// 3. InventoryModify 修改库存
		err = tx.Model(&models.InventoryModify{}).
			Where("goods_id = ?", goodID).Unscoped().Delete(&models.InventoryModify{}).Error
		if err != nil {
			return err
		}

		//        // 4. InventoryHistory 庫存記錄
		//        err = tx.Model(&models.InventoryHistory{}).
		//            Where("goods_id = ?", goodID).Delete(&models.InventoryModify{}).Error
		//        if err != nil {
		//            return err
		//        }

		// 5. InventoryHistoryDetailed 庫存記錄詳細
		err = tx.Model(&models.InventoryHistoryDetailed{}).
			Where("goods_id = ?", goodID).Unscoped().Delete(&models.InventoryModify{}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Simple) Search(keyword string, account string) ([]models.TemporaryGoodsInventories, error) {
	var goods []models.TemporaryGoodsInventories

	if keyword != "" {
		err := s.DB().Raw(`select g.*,i.quantity ,i."cost" as total_cost from goods g left join inventories i on g.barcode  = i.barcode and i.account  = ? where (g.name like ? or g.barcode like ?) order by g.created_at desc limit 30`, account, strings.ReplaceAll(`%P%`, "P", keyword), strings.ReplaceAll(`%P%`, "P", keyword)).Order("g.created_at").Limit(30).Scan(&goods).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := s.DB().Raw(`select g.*,i.quantity ,i."cost" as total_cost from goods g left join inventories i on g.barcode  = i.barcode and i.account  = ? order by g.created_at desc limit 30`, account).Scan(&goods).Error
		if err != nil {
			return nil, err
		}
	}

	for i, k := range goods {
		if k.TotalCost > 0 {
			cost := k.TotalCost / float64(k.Quantity)
			goods[i].Cost = cost
		}
	}

	return goods, nil
}

func (s *Simple) Good(barcodes string, account string) (*models.TemporaryGoodsInventories, error) {
	var good models.Goods
	err := s.DB().Model(&models.Goods{}).Where("barcode = ?", barcodes).First(&good).Error
	if err != nil {
		return nil, err
	}

	var inv models.Inventory
	err = s.DB().Model(&models.Inventory{}).Where("goods_id = ?", good.ID).Where("account = ?", account).First(&inv).Error
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			err = s.DB().Model(&models.Inventory{}).Create(&models.Inventory{
				BasicModel: models.BasicModel{ID: xid.New().String()},
				GoodsID:    good.ID,
				Barcode:    barcodes,
				Account:    account,
			}).Error
			if err != nil {
				return nil, err
			}
		} else {
			log.Println(err)
			return nil, err
		}
	}

	return &models.TemporaryGoodsInventories{
		Goods:     good,
		Quantity:  inv.Quantity,
		TotalCost: inv.Cost,
	}, nil

	//var good models.TemporaryGoodsInventories
	//err := s.DB().Raw(`select g.*,i.quantity ,i."cost" from goods g left join inventories i on g.barcode  = i.barcode and i.account  = ?  and g.barcode = ? `, account, barcodes).Scan(&good).Error
	//if err != nil {
	//	return nil, err
	//}
}
