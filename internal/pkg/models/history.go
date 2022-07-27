package models

type HistoryType string

const (
	HistoryTypeInventoryModify HistoryType = "inventory_modify"
	HistoryTypeInventory       HistoryType = "inventory"
)

type History struct {
	BasicModel
	Account     string      `gorm:"type:varchar(300);index" json:"account"` // 創建用戶
	HistoryType HistoryType `gorm:"type:varchar(300);index" json:"history_type"`
}
