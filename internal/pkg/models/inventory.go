package models

// Inventory 庫存
type Inventory struct {
	BasicModel
	GoodsID string `gorm:"type:varchar(300);index" json:"goods_id"` // 商品
	Account string `gorm:"type:varchar(300);index" json:"account"`  // 創建用戶
}

// InventoryHistory 庫存記錄
type InventoryHistory struct {
	BasicModel
	GoodsID        string  `gorm:"type:varchar(300);index" json:"goods_id"`   // 商品
	Account        string  `gorm:"type:varchar(300);index" json:"account"`    // 創建用戶
	TotalPrice     float64 `gorm:"type:NUMERIC(20,7)" json:"total_price"`     // 總價
	TotalCost      float64 `gorm:"type:NUMERIC(20,7)" json:"total_cost"`      // 縂成本
	GrossProfit    float64 `gorm:"type:NUMERIC(20,7)" json:"gross_profit"`    // 毛利
	NumberProducts int     `gorm:"type:NUMERIC(20,7)" json:"number_products"` // 商品數量
}

// InventoryHistoryDetailed 庫存記錄詳細
type InventoryHistoryDetailed struct {
	BasicModel
	OrderID        string  `gorm:"type:varchar(300);index" json:"order_id"`   // order id
	GoodsID        string  `gorm:"type:varchar(300);index" json:"goods_id"`   // 商品
	Account        string  `gorm:"type:varchar(300);index" json:"account"`    // 創建用戶
	TotalPrice     float64 `gorm:"type:NUMERIC(20,7)" json:"total_price"`     // 總價
	TotalCost      float64 `gorm:"type:NUMERIC(20,7)" json:"total_cost"`      // 縂成本
	GrossProfit    float64 `gorm:"type:NUMERIC(20,7)" json:"gross_profit"`    // 毛利
	NumberProducts int     `gorm:"type:NUMERIC(20,7)" json:"number_products"` // 商品數量
}
