package request

type AddGoods struct {
	Barcode string  `json:"barcode"  binding:"required"`
	Name    string  `json:"name"  binding:"required"`
	Spec    string  `json:"spec"  `
	Cost    float64 `json:"cost"  binding:"required"`
	Price   float64 `json:"price"  binding:"required"`
	Brand   string  `json:"brand"  `
	MadeIn  string  `json:"made_in" `
	Img     string  `json:"img" `
}

type ModifyGoods struct {
	ID     string  `json:"id"  binding:"required"`
	Name   string  `json:"name"  binding:"required"`
	Spec   string  `json:"spec"  `
	Cost   float64 `json:"cost"  binding:"required"`
	Price  float64 `json:"price"  binding:"required"`
	Brand  string  `json:"brand"  `
	MadeIn string  `json:"made_in" `
	Img    string  `json:"img" `
}

type CoreExport struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}
