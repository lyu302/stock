package model

import "time"

type Quote struct {
	IDModel
	StockID   int       `gorm:"column:stock_id; index"`
	Open      float64   `gorm:"column:open"`
	High      float64   `gorm:"column:high"`
	Low       float64   `gorm:"column:low"`
	Close     float64   `gorm:"column:close"`
	Volume    int64   `gorm:"column:volume"`
	Date      time.Time `gorm:"column:date; type:date"`
	ChangePercent float64  `gorm:"column:change_percent"`

	Stock     *Stock
}

func (q *Quote) TableName() string {
	return "quotes"
}
