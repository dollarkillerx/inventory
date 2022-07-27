package simple

import "github.com/dollarkillerx/inventory/internal/pkg/models"

func (s *Simple) Goods() {

}

func (s *Simple) Good(barcodes string, account string) (*models.TemporaryGoodsInventories, error) {
	var good models.TemporaryGoodsInventories
	err := s.DB().Raw(`select g.*,i.quantity ,i."cost" from goods g left join inventories i on g.barcode  = i.barcode where i.account  = ?  and g.barcode = ? `, account, barcodes).Scan(&good).Error
	if err != nil {
		return nil, err
	}

	return &good, nil
}
