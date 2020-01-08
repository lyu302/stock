package db

import (
	"github.com/lyu302/stock/db/model"
	"time"
)

type Dao interface {
	AddModel(model.Interface) error
	UpdateModel(model.Interface) error
}

type DelDao interface {
	DeleteModel(ID string, arg ...interface{}) error
}

//StockDao: stock dao interface
type StockDao interface {
	Dao
	FindStockBySymbolCode(symbol, code string) (*model.Stock, error)
	ListAllStocks() ([]*model.Stock, error)
}

type QuoteDao interface {
	Dao
	ListQuotesByStockID(stockID int) ([]*model.Quote, error)
	ListQuotesByStockIDAndTime(stockID int, start, end time.Time) ([]*model.Quote, error)
	ListQuotesByZeroPercent() ([]*model.Quote, error)
	FindQuoteByLastDate(stockID int, date time.Time) (*model.Quote, error)
	ListQuotesByTime(start, end time.Time) ([]*model.Quote, error)
}