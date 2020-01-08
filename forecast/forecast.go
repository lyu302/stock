package forecast

import "github.com/lyu302/stock/db/model"

type Forecast interface {
}

type Result struct {
	Stock     *model.Stock
	Hit       [][]*model.Quote
	Before    [][]*model.Quote
	After     [][]*model.Quote
}