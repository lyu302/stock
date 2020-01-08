package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/lyu302/stock/db/model"
	"log"
)

type StockDaoImpl struct {
	DB *gorm.DB
}

func (s *StockDaoImpl) AddModel(m model.Interface) error {
	stock := m.(*model.Stock)
	var oldStock model.Stock

	if ok := s.DB.Where("symbol = ? and code = ?", stock.Symbol, stock.Code).Find(&oldStock).RecordNotFound(); ok {
		if err := s.DB.Create(stock).Error; err != nil {
			return err
		}
	} else {
		log.Printf("Symbol %s And Code %s Has Exist In Stock", stock.Symbol, stock.Code)
	}

	return nil
}

func (s *StockDaoImpl) UpdateModel(m model.Interface) error  {
	stock := m.(*model.Stock)
	if err := s.DB.Save(stock).Error; err != nil {
		return err
	}

	return nil
}

func (s *StockDaoImpl) FindStockBySymbolCode(symbol, code string) (*model.Stock, error)  {
	var stock model.Stock
	if err := s.DB.Where("symbol = ? and code = ?", symbol, code).Find(&stock).Error; err != nil {
		return nil, err
	}

	return &stock, nil
}

func (s *StockDaoImpl) ListAllStocks() ([]*model.Stock, error) {
	var stocks []*model.Stock
	if err := s.DB.Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}