package models

type StatisticsBand struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

//type Statistics struct {
//	TotalInventoryQuantity int     `json:"total_inventory_quantity"` // 縂庫存數量
//	TotalInventoryCost     float64 `json:"total_inventory_cost"`     // 縂庫存成本
//	Turnover24             float64 `json:"turnover24"`               // 24h 流水
//	Profit24               float64 `json:"profit24"`                 // 24h 利潤
//	TurnoverWeek           float64 `json:"turnover_week"`            // 一周 流水
//	ProfitWeek             float64 `json:"profit_week"`              // 一周 利潤
//	TurnoverMoon           float64 `json:"turnover_moon"`            // 月 流水
//	ProfitMoon             float64 `json:"profit_moon"`              // 月 利潤
//}
