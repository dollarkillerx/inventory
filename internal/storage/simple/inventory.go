package simple

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"

	"log"
	"strings"
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
	err = begin.Model(&models.InventoryHistory{}).Create(&models.InventoryHistory{
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

	err = begin.Model(&models.InventoryHistoryDetailed{}).Create(&models.InventoryHistoryDetailed{
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
	err = begin.Model(&models.Inventory{}).
		Where("goods_id = ?", goodsId).
		Where("account = ?", account).First(&inv).Error
	if err != nil {
		return err
	}

	inv.Quantity += numberProducts
	inv.Cost += cost

	err = begin.Model(&models.Inventory{}).
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

	var inv models.Inventory
	err = begin.Model(&models.Inventory{}).
		Where("goods_id = ?", goodsId).
		Where("account = ?", account).First(&inv).Error
	if err != nil {
		return err
	}

	if inv.Cost > 0 && inv.Quantity > 0 {
		cst := inv.Cost / float64(inv.Quantity) // 計算成本單價
		cost = cst * float64(numberProducts)    // 計算成本
	}

	ihid := xid.New().String()
	err = begin.Model(&models.InventoryHistory{}).Create(&models.InventoryHistory{
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

	err = begin.Model(&models.InventoryHistoryDetailed{}).Create(&models.InventoryHistoryDetailed{
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

	inv.Quantity -= numberProducts
	inv.Cost -= cost

	err = begin.Model(&models.Inventory{}).
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

func (s *Simple) IORevoke(orderID string, account string) (result []models.InventoryHistoryDetailed, err error) {
	s.inventoryMu.Lock()
	defer func() {
		s.inventoryMu.Unlock()
	}()

	begin := s.db.Begin()

	defer func() {
		if err == nil {
			begin.Commit()
		} else {
			begin.Rollback()
		}
	}()

	var orderDetailed models.InventoryHistoryDetailed
	err = begin.Model(&models.InventoryHistoryDetailed{}).
		Where("id = ?", orderID).
		Where("account = ?", account).First(&orderDetailed).Error
	if err != nil {
		return nil, err
	}

	var count int64
	err = begin.Model(&models.InventoryHistoryDetailed{}).Where("order_id = ?", orderDetailed.OrderID).Count(&count).Error
	if err != nil {
		return nil, err
	}

	if count != 1 {
		return nil, errors.New("當前訂單存在子訂單 無權刪除")
	}

	var goodsInventories models.Inventory
	err = begin.Model(&models.Inventory{}).
		Where("account = ?", account).
		Where("goods_id = ?", orderDetailed.GoodsID).First(&goodsInventories).Error
	if err != nil {
		return nil, err
	}

	err = begin.Model(&models.InventoryHistoryDetailed{}).Where("id = ?", orderID).Delete(&models.InventoryHistoryDetailed{}).Error
	if err != nil {
		return nil, err
	}

	err = begin.Model(&models.InventoryHistory{}).Where("id = ?", orderDetailed.OrderID).Delete(&models.InventoryHistory{}).Error
	if err != nil {
		return nil, err
	}

	switch orderDetailed.InventoryType {
	case models.InventoryHistoryTypeWarehousing:
		err = begin.Model(&models.Inventory{}).
			Where("account = ?", account).
			Where("goods_id = ?", orderDetailed.GoodsID).
			Updates(map[string]interface{}{
				"quantity": goodsInventories.Quantity - orderDetailed.NumberProducts,
				"cost":     goodsInventories.Cost - orderDetailed.TotalCost,
			}).Error
		if err != nil {
			return nil, err
		}
	case models.InventoryHistoryTypeDepot:
		err = begin.Model(&models.Inventory{}).
			Where("account = ?", account).
			Where("goods_id = ?", orderDetailed.GoodsID).
			Updates(map[string]interface{}{
				"quantity": goodsInventories.Quantity + orderDetailed.NumberProducts,
				"cost":     goodsInventories.Cost + orderDetailed.TotalCost,
			}).Error
		if err != nil {
			return nil, err
		}
	}

	err = begin.Model(&models.InventoryHistoryDetailed{}).
		Where("goods_id = ?", orderDetailed.GoodsID).
		Where("account = ?", account).Order("created_at desc").Find(&result).Error

	return
}

func (s *Simple) ResetStatistics() (err error) {
	s.inventoryMu.Lock()
	defer func() {
		s.inventoryMu.Unlock()
	}()

	begin := s.db.Begin()

	defer func() {
		if err == nil {
			begin.Commit()
		} else {
			begin.Rollback()
		}
	}()

	//err = s.resetStatisticsInternal(begin)
	//if err != nil {
	//	return err
	//}

	var ih []models.InventoryHistory
	err = begin.Model(&models.InventoryHistory{}).Find(&ih).Error
	if err != nil {
		return err
	}

	for _, v := range ih {
		var ihd []models.InventoryHistoryDetailed
		err = begin.Model(&models.InventoryHistoryDetailed{}).Where("order_id = ?", v.ID).Find(&ihd).Error
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				err = nil
			} else {
				log.Println(err)
				return err
			}
		}

		var totalPrice float64  // 總價
		var totalCost float64   // 縂成本
		var grossProfit float64 // 毛利
		var numberProducts int  // 商品數量

		for _, vc := range ihd {
			totalCost += vc.TotalCost
			totalPrice += vc.TotalPrice
			grossProfit += vc.GrossProfit
			numberProducts += vc.NumberProducts
		}

		err = begin.Model(&models.InventoryHistory{}).
			Where("id = ?", v.ID).Updates(map[string]interface{}{
			"total_price":     totalPrice,
			"total_cost":      totalCost,
			"gross_profit":    grossProfit,
			"number_products": numberProducts,
		}).Error
		if err != nil {
			return err
		}
	}

	var ins []models.Inventory
	err = begin.Model(&models.Inventory{}).Find(&ins).Error
	if err != nil {
		return err
	}

	for _, v := range ins {
		var ih []models.InventoryHistoryDetailed
		err = begin.Model(&models.InventoryHistoryDetailed{}).
			Where("account = ?", v.Account).Where("goods_id = ?", v.GoodsID).Find(&ih).Error
		if err != nil {
			return err
		}

		var quantity int // 库存数量
		var cost float64 // 总成本

		for _, vc := range ih {
			switch vc.InventoryType {
			case models.InventoryHistoryTypeWarehousing:
				quantity += vc.NumberProducts
				cost += vc.TotalCost
			case models.InventoryHistoryTypeDepot:
				quantity -= vc.NumberProducts
				cost -= vc.TotalCost
			}
		}

		err = begin.Model(&models.Inventory{}).Where("id = ?", v.ID).Updates(map[string]interface{}{
			"quantity": quantity,
			"cost":     cost,
		}).Error
		if err != nil {
			return err
		}
	}

	return
}

func (s *Simple) resetStatisticsInternal(db *gorm.DB) (err error) {
	var ihd []models.InventoryHistoryDetailed
	err = db.Model(&models.InventoryHistoryDetailed{}).Find(&ihd).Error
	if err != nil {
		return err
	}

	for _, v := range ihd {
		var goods models.Goods
		err = db.Model(&models.Goods{}).Where("id = ?", v.GoodsID).First(&goods).Error
		if err != nil {
			return err
		}

		var totalPrice float64  // 總價
		var totalCost float64   // 縂成本
		var grossProfit float64 // 毛利

		totalPrice = goods.Price * float64(v.NumberProducts)
		totalCost = goods.Cost * float64(v.NumberProducts)
		grossProfit = totalPrice - totalCost

		err = db.Model(&models.InventoryHistoryDetailed{}).
			Where("id = ?", v.ID).Updates(map[string]interface{}{
			"total_price":  totalPrice,
			"total_cost":   totalCost,
			"gross_profit": grossProfit,
		}).Error
		if err != nil {
			return err
		}
	}
	return
}

func (s *Simple) IOList(account string, limit int, offset int) (count int64, ids []models.InventoryHistoryDetailed, err error) {
	err = s.db.Model(&models.InventoryHistoryDetailed{}).
		Where("account = ?", account).
		Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	if limit <= 0 || limit > 20 {
		limit = 20
	}

	err = s.db.Model(&models.InventoryHistoryDetailed{}).
		Where("account = ?", account).
		Limit(limit).Offset(offset).
		Order("created_at desc").Find(&ids).Error

	var cmap = map[string]struct{}{}
	for _, v := range ids {
		cmap[v.GoodsID] = struct{}{}
	}

	var goodsList []string
	for k := range cmap {
		goodsList = append(goodsList, k)
	}

	var goods []models.Goods
	err = s.db.Model(&models.Goods{}).Where("id in ?", goodsList).Find(&goods).Error
	if err != nil {
		return 0, nil, err
	}

	for i, v := range ids {
		for _, vc := range goods {
			if v.GoodsID == vc.ID {
				ids[i].GoodsName = vc.Name
			}
		}
	}

	return
}
