package request

type Warehousing struct {
	Barcode        string  `json:"barcode" binding:"required"` // 条码
	Cost           float64 `json:"cost" binding:"required"`    // 成本
	NumberProducts int     `json:"number_products" binding:"required"`
	Remark         string  `json:"remark"`
}

type OutStock struct {
	Barcode        string  `json:"barcode" binding:"required"` // 条码
	Cost           float64 `json:"cost" binding:"required"`    // 成本
	Price          float64 `json:"price" binding:"required"`   // 总价
	NumberProducts int     `json:"number_products" binding:"required"`
	Remark         string  `json:"remark"`
}

type IORevoke struct {
	OrderID string `json:"order_id" binding:"required"`
}
