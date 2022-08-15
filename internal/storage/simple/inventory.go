package simple

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/rs/xid"
)

func (s *Simple) WareHousing(goodsId string, barcode string, account string, cost float64, numberProducts int, remark string) (err error) {
	s.inventoryMu.Lock()
	defer func() {
		s.inventoryMu.Unlock()
	}()

	begin := s.db.Begin()
	defer func() {
		if err == nil {
			_ = begin.Commit().Error
		} else {
			_ = begin.Rollback().Error
		}
	}()

	ihid := xid.New().String()
	err = s.DB().Model(&models.InventoryHistory{}).Create(&models.InventoryHistory{
		BasicModel:     models.BasicModel{ID: ihid},
		InventoryType:  models.InventoryHistoryTypeWarehousing,
		Account:        account,
		TotalCost:      cost,
		NumberProducts: numberProducts,
		Remark:         remark,
	}).Error
	if err != nil {
		return err
	}

	err = s.DB().Model(&models.InventoryHistoryDetailed{}).Create(&models.InventoryHistoryDetailed{
		BasicModel:     models.BasicModel{ID: xid.New().String()},
		InventoryType:  models.InventoryHistoryTypeWarehousing,
		OrderID:        ihid,
		Barcode:        barcode,
		GoodsID:        goodsId,
		Account:        account,
		TotalCost:      cost,
		NumberProducts: numberProducts,
		Remark:         remark,
	}).Error
	if err != nil {
		return err
	}

	var inv models.Inventory
	err = s.DB().Model(&models.Inventory{}).
		Where("goods_id = ?", goodsId).
		Where("account = ?", account).First(&inv).Error
	if err != nil {
		return err
	}

	inv.Quantity += numberProducts
	inv.Cost += cost

	err = s.DB().Model(&models.Inventory{}).
		Where("goods_id = ?", goodsId).
		Where("account = ?", account).
		Updates(map[string]interface{}{
			"quantity": inv.Quantity,
			"cost":     inv.Cost,
		}).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Simple) OutStock(goodsId string, barcode string, account string, cost float64, numberProducts int, price float64, remark string) (err error) {
	s.inventoryMu.Lock()
	defer func() {
		s.inventoryMu.Unlock()
	}()

	begin := s.db.Begin()
	defer func() {
		if err == nil {
			_ = begin.Commit().Error
		} else {
			_ = begin.Rollback().Error
		}
	}()

	ihid := xid.New().String()
	err = s.DB().Model(&models.InventoryHistory{}).Create(&models.InventoryHistory{
		BasicModel:     models.BasicModel{ID: ihid},
		InventoryType:  models.InventoryHistoryTypeDepot,
		Account:        account,
		TotalCost:      cost,
		TotalPrice:     price,
		GrossProfit:    price - cost,
		NumberProducts: numberProducts,
		Remark:         remark,
	}).Error
	if err != nil {
		return err
	}

	err = s.DB().Model(&models.InventoryHistoryDetailed{}).Create(&models.InventoryHistoryDetailed{
		BasicModel:     models.BasicModel{ID: xid.New().String()},
		InventoryType:  models.InventoryHistoryTypeDepot,
		OrderID:        ihid,
		Barcode:        barcode,
		GoodsID:        goodsId,
		Account:        account,
		TotalCost:      cost,
		TotalPrice:     price,
		GrossProfit:    price - cost,
		NumberProducts: numberProducts,
		Remark:         remark,
	}).Error
	if err != nil {
		return err
	}

	var inv models.Inventory
	err = s.DB().Model(&models.Inventory{}).
		Where("goods_id = ?", goodsId).
		Where("account = ?", account).First(&inv).Error
	if err != nil {
		return err
	}

	inv.Quantity -= numberProducts
	inv.Cost -= cost

	err = s.DB().Model(&models.Inventory{}).
		Where("goods_id = ?", goodsId).
		Where("account = ?", account).
		Updates(map[string]interface{}{
			"quantity": inv.Quantity,
			"cost":     inv.Cost,
		}).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Simple) IOHistory(goodsID string, account string) (resp []models.InventoryHistoryDetailed, err error) {
	err = s.db.Model(&models.InventoryHistoryDetailed{}).
		Where("goods_id = ?", goodsID).
		Where("account = ?", account).Order("created_at desc").Find(&resp).Error

	return
}
