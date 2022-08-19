package models

// Inventory 庫存
type Inventory struct {
	BasicModel
	GoodsID  string  `gorm:"type:varchar(300);index" json:"goods_id"` // 商品
	Barcode  string  `gorm:"type:varchar(300);index" json:"barcode"`  // 條形碼
	Account  string  `gorm:"type:varchar(300);index" json:"account"`  // 創建用戶
	Quantity int     `json:"quantity"`                                // 库存数量
	Cost     float64 `gorm:"type:NUMERIC(20,7)" json:"cost"`          // 总成本
}

// InventoryModify 修改库存
type InventoryModify struct {
	BasicModel
	GoodsID        string  `gorm:"type:varchar(300);index" json:"goods_id"`        // 商品
	Barcode        string  `gorm:"type:varchar(300);index" json:"barcode"`         // 條形碼
	Account        string  `gorm:"type:varchar(300);index" json:"account"`         // 創建用戶
	BeforeQuantity int     `gorm:"type:varchar(300);index" json:"before_quantity"` // 库存数量
	BeforeCost     float64 `gorm:"type:NUMERIC(20,7)" json:"before_cost"`          // 成本
	Quantity       int     `gorm:"type:varchar(300);index" json:"quantity"`        // 库存数量
	Cost           float64 `gorm:"type:NUMERIC(20,7)" json:"cost"`                 // 成本
}

type InventoryHistoryType string

const (
	InventoryHistoryTypeWarehousing InventoryHistoryType = "warehousing" // 入库
	InventoryHistoryTypeDepot       InventoryHistoryType = "depot"       // 出库
)

// InventoryHistory 庫存記錄
type InventoryHistory struct {
	BasicModel
	InventoryType  InventoryHistoryType `gorm:"type:varchar(300);index" json:"inventory_type"`
	Account        string               `gorm:"type:varchar(300);index" json:"account"`    // 創建用戶
	TotalPrice     float64              `gorm:"type:NUMERIC(20,7)" json:"total_price"`     // 總價
	TotalCost      float64              `gorm:"type:NUMERIC(20,7)" json:"total_cost"`      // 縂成本
	GrossProfit    float64              `gorm:"type:NUMERIC(20,7)" json:"gross_profit"`    // 毛利
	NumberProducts int                  `gorm:"type:NUMERIC(20,7)" json:"number_products"` // 商品數量
	Remark         string               `gorm:"type:varchar(700)" json:"remark"`
}

// InventoryHistoryDetailed 庫存記錄詳細
type InventoryHistoryDetailed struct {
	BasicModel
	InventoryType  InventoryHistoryType `gorm:"type:varchar(300);index" json:"inventory_type"`
	OrderID        string               `gorm:"type:varchar(300);index" json:"order_id"` // order id
	Barcode        string               `gorm:"type:varchar(300);index" json:"barcode"`  // 條形碼
	GoodsID        string               `gorm:"type:varchar(300);index" json:"goods_id"` // 商品
	Account        string               `gorm:"type:varchar(300);index" json:"account"`  // 創建用戶
	TotalPrice     float64              `gorm:"type:NUMERIC(20,7)" json:"total_price"`   // 總價
	GrossProfit    float64              `gorm:"type:NUMERIC(20,7)" json:"gross_profit"`  // 毛利
	TotalCost      float64              `gorm:"type:NUMERIC(20,7)" json:"total_cost"`    // 縂成本
	NumberProducts int                  `gorm:"type:int(8)" json:"number_products"`      // 商品數量
	Remark         string               `gorm:"type:varchar(700)" json:"remark"`
}
