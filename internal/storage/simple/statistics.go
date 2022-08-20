package simple

import (
	"fmt"
	"time"

	"github.com/dollarkillerx/inventory/internal/pkg/models"
)

func (s *Simple) Statistics(account string) (scs []models.StatisticsBand, err error) {
	var ins []models.Inventory
	err = s.db.Model(&models.Inventory{}).Where("account = ?", account).Find(&ins).Error
	if err != nil {
		return nil, err
	}

	var totalCost float64
	var totalQuantity int
	for _, v := range ins {
		totalCost += v.Cost
		totalQuantity += v.Quantity
	}

	scs = append(scs, models.StatisticsBand{
		Key:   "庫存成本",
		Value: []string{fmt.Sprintf("%.2f", totalCost)},
	})

	scs = append(scs, models.StatisticsBand{
		Key:   "庫存數量",
		Value: []string{fmt.Sprintf("%d", totalQuantity)},
	})

	type timeFilter struct {
		Key   string    `json:"key"`
		Start time.Time `json:"time"`
		End   time.Time `json:"end"`
	}

	t := time.Now()
	ts := time.Now().AddDate(0, 0, -1)
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	moon := t.AddDate(0, 0, -t.Day()+1)
	lastMonthFirstDay := t.AddDate(0, -1, -t.Day()+1)
	lastMonthStart := time.Date(lastMonthFirstDay.Year(), lastMonthFirstDay.Month(), lastMonthFirstDay.Day(), 0, 0, 0, 0, t.Location())

	timeFilters := []timeFilter{
		{
			Key:   "今天",
			Start: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
			End:   time.Now().Add(6 * time.Hour),
		},
		{
			Key:   "昨天",
			Start: time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location()),
			End:   time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
		},
		{
			Key:   "本周",
			Start: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset),
			End:   time.Now().Add(6 * time.Hour),
		},
		{
			Key:   "上周",
			Start: GetLastWeekFirstDate(),
			End:   GetFirstDateOfWeek(),
		},
		{
			Key:   "本月",
			Start: time.Date(moon.Year(), moon.Month(), moon.Day(), 0, 0, 0, 0, moon.Location()),
			End:   time.Now().Add(6 * time.Hour),
		},
		{
			Key:   "上個月",
			Start: lastMonthStart,
			End:   time.Date(moon.Year(), moon.Month(), moon.Day(), 0, 0, 0, 0, moon.Location()),
		},
	}

	for _, v := range timeFilters {
		var ihs []models.InventoryHistory
		err = s.db.Model(&models.InventoryHistory{}).
			Where("account = ?", account).
			Where("inventory_type = ?", models.InventoryHistoryTypeDepot).
			Where("created_at >= ?", v.Start).
			Where("created_at < ?", v.End).Find(&ihs).Error
		if err != nil {
			return nil, err
		}

		var turnover float64 // 流水
		var profit float64   // 利潤
		for _, vc := range ihs {
			turnover += vc.TotalPrice
			profit += vc.GrossProfit
		}

		scs = append(scs, models.StatisticsBand{
			Key: v.Key,
			Value: []string{
				fmt.Sprintf("流水: %.2f", turnover),
				fmt.Sprintf("利潤: %.2f", profit),
			},
		})
	}

	return
}

/**
获取本周周一的日期
*/
func GetFirstDateOfWeek() (weekMonday time.Time) {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStartDate
}

/**
获取上周的周一日期
*/
func GetLastWeekFirstDate() (weekMonday time.Time) {
	thisWeekMonday := GetFirstDateOfWeek()
	lastWeekMonday := thisWeekMonday.AddDate(0, 0, -7)
	return lastWeekMonday
}
