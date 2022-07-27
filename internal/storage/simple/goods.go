package simple

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/rs/xid"
)

func (s *Simple) Goods() {

}

func (s *Simple) Good(barcodes string, account string) (*models.TemporaryGoodsInventories, error) {
	var good models.Goods
	err := s.DB().Model(&models.Goods{}).Where("barcode = ?", barcodes).First(&good).Error
	if err != nil {
		return nil, err
	}

	var inv models.Inventory
	err = s.DB().Model(&models.Inventory{}).Where("barcode = ?", barcodes).Where("account = ?", account).First(&inv).Error
	if err != nil {
		err := s.DB().Model(&models.Inventory{}).Create(&models.Inventory{
			BasicModel: models.BasicModel{ID: xid.New().String()},
			GoodsID:    good.ID,
			Barcode:    barcodes,
			Account:    account,
		}).Error
		if err != nil {
			return nil, err
		}
	}

	return &models.TemporaryGoodsInventories{
		Goods:    good,
		Quantity: inv.Quantity,
		Cost:     inv.Cost,
	}, nil

	//var good models.TemporaryGoodsInventories
	//err := s.DB().Raw(`select g.*,i.quantity ,i."cost" from goods g left join inventories i on g.barcode  = i.barcode and i.account  = ?  and g.barcode = ? `, account, barcodes).Scan(&good).Error
	//if err != nil {
	//	return nil, err
	//}
}
