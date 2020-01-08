package model

import "time"

type IDModel struct {
	ID uint `gorm:"column:ID;primary_key"`
}

type TimeModel struct {
	CreatedAt time.Time `gorm:"column:created_at; default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Interface interface {
	TableName() string
}