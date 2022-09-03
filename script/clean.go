package main

import (
	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/dollarkillerx/inventory/internal/storage/simple"
	"gorm.io/gorm"

	"log"
)

func main() {
	newSimple, err := simple.NewSimple(&conf.CONF.PgSQLConfig)
	if err != nil {
		panic(err)
	}

	err = clean(newSimple.DB())
	if err != nil {
		log.Fatalln(err)
	}
}

func clean(db *gorm.DB) (err error) {
	begin := db.Begin()
	defer func() {
		if err == nil {
			begin.Commit()
		} else {
			begin.Rollback()
		}
	}()

	var ins []models.Inventory
	err = begin.Model(&models.Inventory{}).Find(&ins).Error
	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range ins {
		var idd []models.InventoryHistoryDetailed
		err = begin.Model(&models.InventoryHistoryDetailed{}).
			Where("goods_id = ?", v.GoodsID).Where("account = ?", v.Account).Find(&idd).Error
		if err != nil {
			log.Println(err)
			return err
		}

		var total int
		var cost float64

		for _, vc := range idd {
			switch vc.InventoryType {
			case models.InventoryHistoryTypeWarehousing:
				total += vc.NumberProducts
				cost += vc.TotalCost
			case models.InventoryHistoryTypeDepot:
				total -= vc.NumberProducts
				cost -= vc.TotalCost
			}
		}

		if total == 0 {
			cost = 0
		}

		err = begin.Model(&models.Inventory{}).
			Where("account = ?", v.Account).
			Where("id = ?", v.ID).Updates(map[string]interface{}{
			"quantity": total,
			"cost":     cost,
		}).Error
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
