package db

import mysqlDao "github.com/lyu302/stock/db/mysql"

//current db is mysql, use mysqlDao
//also support other dialects: postgres sqlite mssql

func (m *Manager) StockDao() StockDao {
	return &mysqlDao.StockDaoImpl {
		DB: m.db,
	}
}

func (m *Manager) QuoteDao() QuoteDao  {
	return &mysqlDao.QuoteDaoImpl {
		DB:m.db,
	}
}
