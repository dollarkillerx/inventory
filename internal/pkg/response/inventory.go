package response

import "github.com/dollarkillerx/inventory/internal/pkg/models"

type IOList struct {
	Count int64                             `json:"count"`
	Items []models.InventoryHistoryDetailed `json:"items"`
}
