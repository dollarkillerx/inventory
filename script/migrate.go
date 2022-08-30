package main

import (
	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/dollarkillerx/inventory/internal/storage/simple"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"log"
)

func main() {
	newSimple, err := simple.NewSimple(&conf.CONF.PgSQLConfig)
	if err != nil {
		panic(err)
	}

	err = migrate(newSimple.DB())
	if err != nil {
		log.Fatalln(err)
	}
}

func migrate(db *gorm.DB) (err error) {
	begin := db.Begin()
	defer func() {
		if err == nil {
			begin.Commit()
		} else {
			begin.Rollback()
		}
	}()

	account := "991"
	notAccount := "119"

	var ivths []models.InventoryHistoryDetailed
	err = begin.Model(&models.InventoryHistoryDetailed{}).Where("account != ?", notAccount).Find(&ivths).Error
	if err != nil {
		log.Println(err)
		return err
	}

	for i := range ivths {
		ivths[i].ID = uuid.New().String()
		ivths[i].Account = account
		ivths[i].OrderID = uuid.New().String()
	}

	var ihs []models.InventoryHistory
	var ityMap = map[string]models.Inventory{}
	for _, v := range ivths {
		inventory, ok := ityMap[v.GoodsID]
		if !ok {
			inventory.BasicModel = models.BasicModel{ID: xid.New().String()}
			inventory.Account = account
			inventory.Barcode = v.Barcode
			inventory.GoodsID = v.GoodsID
		}
		if v.InventoryType == models.InventoryHistoryTypeDepot {
			inventory.Quantity -= v.NumberProducts
			inventory.Cost -= v.TotalCost
		} else {
			inventory.Quantity += v.NumberProducts
			inventory.Cost += v.TotalCost
		}
		ityMap[v.GoodsID] = inventory

		ihs = append(ihs, models.InventoryHistory{
			BasicModel: models.BasicModel{
				ID:        v.OrderID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
				DeletedAt: v.DeletedAt,
			},
			InventoryType:  v.InventoryType,
			Account:        v.Account,
			TotalPrice:     v.TotalPrice,
			TotalCost:      v.TotalCost,
			GrossProfit:    v.GrossProfit,
			NumberProducts: v.NumberProducts,
			Remark:         v.Remark,
		})
	}

	var ity []models.Inventory
	for ic := range ityMap {
		ity = append(ity, ityMap[ic])
	}

	err = begin.Model(&models.Inventory{}).Create(&ity).Error
	if err != nil {
		log.Println(err)
		return err
	}

	err = begin.Model(&models.InventoryHistory{}).Create(&ihs).Error
	if err != nil {
		log.Println(err)
		return err
	}

	err = begin.Model(&models.InventoryHistoryDetailed{}).Create(&ivths).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func migrate2(db *gorm.DB) (err error) {
	begin := db.Begin()
	defer func() {
		if err == nil {
			begin.Commit()
		} else {
			begin.Rollback()
		}
	}()

	account := "991"
	notAccount := "119"

	var invs []models.Inventory
	err = begin.Model(&models.Inventory{}).Where("account != ?", notAccount).Find(&invs).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var ivi []models.InventoryHistory
	err = begin.Model(&models.InventoryHistory{}).Where("account != ?", notAccount).Find(&ivi).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var ivths []models.InventoryHistoryDetailed
	err = begin.Model(&models.InventoryHistoryDetailed{}).Where("account != ?", notAccount).Find(&ivths).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var newInvs []models.Inventory
	var vm = map[string][]models.Inventory{}
	for _, v := range invs {
		vm[v.GoodsID] = append(vm[v.GoodsID], v)
	}

	for _, v := range vm {
		var newI = models.Inventory{
			BasicModel: models.BasicModel{ID: xid.New().String()},
		}

		for _, vk := range v {
			newI.GoodsID = vk.GoodsID
			newI.Barcode = vk.Barcode
			newI.Account = account
			newI.Quantity += vk.Quantity
			newI.Cost += vk.Cost
		}

		newInvs = append(newInvs, newI)
	}

	err = begin.Model(&models.Inventory{}).Create(&newInvs).Error
	if err != nil {
		log.Println(err)
		return err
	}

	for i := range ivi {
		ivi[i].ID = xid.New().String()
		ivi[i].Account = account
	}
	err = begin.Model(&models.InventoryHistory{}).Create(&ivi).Error
	if err != nil {
		log.Println(err)
		return err
	}

	for i := range ivths {
		ivths[i].ID = xid.New().String()
		ivths[i].Account = account
	}
	err = begin.Model(&models.InventoryHistoryDetailed{}).Create(&ivths).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
