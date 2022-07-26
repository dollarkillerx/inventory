package models

// Goods 商品
type Goods struct {
	BasicModel
	Barcode   string  `gorm:"type:varchar(300);index" json:"barcode"`    // 條形碼
	Name      string  `gorm:"type:varchar(600);index" json:"name"`       // 商品名稱
	Spec      string  `gorm:"type:varchar(600)" json:"spec"`             // 規格
	Cost      float64 `gorm:"type:NUMERIC(20,7)"  json:"cost"`           // 成本
	Price     float64 `gorm:"type:NUMERIC(20,7)" json:"price"`           // 價格
	Brand     string  `gorm:"type:varchar(300)" json:"brand"`            // 品牌
	MadeIn    string  `gorm:"type:varchar(600)" json:"made_in"`          // 產地
	Img       string  `gorm:"type:varchar(600)" json:"img"`              // img
	ByAccount string  `gorm:"type:varchar(300);index" json:"by_account"` // 創建用戶
}
