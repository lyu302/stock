package model

import "fmt"

type Stock struct {
	IDModel
	Symbol    string  `gorm:"column:symbol; index:idx_symbol_code"`
	Code      string  `gorm:"column:code; index:idx_symbol_code"`
	Name      string  `gorm:"column:name; index"`
	TimeModel
}

func (s *Stock) TableName() string {
	return "stocks"
}

func (s *Stock) SymbolCode() string  {
	return fmt.Sprintf("%s.%s", s.Symbol, s.Code)
}