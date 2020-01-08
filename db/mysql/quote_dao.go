package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/lyu302/stock/db/model"
	"log"
	"time"
)

type QuoteDaoImpl struct {
	DB *gorm.DB
}

func (q *QuoteDaoImpl) AddModel(m model.Interface) error {
	quote := m.(*model.Quote)
	var oldQuote model.Quote

	if ok := q.DB.Where("stock_id = ? and date = ?", quote.StockID, quote.Date).Find(&oldQuote).RecordNotFound(); ok {
		if err := q.DB.Create(quote).Error; err != nil {
			return err
		}
	} else {
		log.Printf("Stock Quote (%s %s) Has Exist In DB", quote.StockID, quote.Date.Format("2006-01-02"))
	}

	return nil
}

func (q *QuoteDaoImpl) UpdateModel(m model.Interface) error {
	quote := m.(*model.Quote)
	if err := q.DB.Save(quote).Error; err != nil {
		return err
	}

	return nil
}

func (q *QuoteDaoImpl) ListQuotesByStockID(stockId int) ([]*model.Quote, error) {
	var quotes []*model.Quote
	if err := q.DB.Where("stock_id = ?", stockId).Find(&quotes).Error; err != nil {
		return nil, err
	}

	return quotes, nil
}

func (q *QuoteDaoImpl) ListQuotesByStockIDAndTime(stockID int, start, end time.Time) ([]*model.Quote, error) {
	var quotes []*model.Quote
	if err := q.DB.Where("stock_id = ? and date between ? and ?", stockID, start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&quotes).Error; err != nil {
		return nil, err
	}

	return quotes, nil
}

// 修复一字涨停、开盘涨停等导致涨幅为0的bug
func (q *QuoteDaoImpl) ListQuotesByZeroPercent() ([]*model.Quote, error) {
	var quotes []*model.Quote
	if err := q.DB.Where("change_percent = 0").Find(&quotes).Error; err != nil {
		return nil, err
	}

	return quotes, nil
}

func (q *QuoteDaoImpl) FindQuoteByLastDate(stockID int, date time.Time) (*model.Quote, error)  {
	var quote model.Quote
	if err := q.DB.Where("stock_id = ? and date < ?", stockID, date).Order("date desc").First(&quote).Error; err != nil {
		return nil, err
	}

	return &quote, nil
}

func (q *QuoteDaoImpl) ListQuotesByTime(start, end time.Time) ([]*model.Quote, error) {
	var quotes []*model.Quote
	if err := q.DB.Preload("Stock").Where("date between ? and ?", start.Format("2006-01-02"), end.Format("2006-01-02")).Find(&quotes).Error; err != nil {
		return nil, err
	}

	return quotes, nil
}