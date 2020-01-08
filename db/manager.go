package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"sync"

	"github.com/lyu302/stock/db/config"
	"github.com/lyu302/stock/db/model"
)

type Manager struct {
	db     *gorm.DB
	config config.Config
	once   sync.Once
	models []model.Interface
}

var DefaultManager *Manager

func NewDefaultManager(config config.Config) (*Manager, error) {
	var err error

	DefaultManager, err = NewManager(config)
	if err != nil {
		return nil, err
	}

	return DefaultManager, nil
}

func NewManager(config config.Config) (*Manager, error)  {
	var (
		db       *gorm.DB
		err      error
		connInfo = config.DbConnectionInfo
	)


	if config.DbType == "mysql" {
		connInfo = fmt.Sprintf("%s?charset=utf8&parseTime=True&loc=Local", connInfo)
	}

	db, err = gorm.Open(config.DbType, connInfo)
	if err != nil {
		return nil, err
	}

	// default set for sql pool
	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(100)

	manager := &Manager{
		db:     db,
		config: config,
		once:   sync.Once{},
	}

	manager.RegisterTableModel()
	manager.CheckTable()
	
	return manager, nil
}

func (m *Manager) Close() error  {
	return m.db.Close()
}

func (m *Manager) Begin() *gorm.DB {
	return m.db.Begin()
}

func (m *Manager) RegisterTableModel() {
	m.models = append(m.models, &model.Stock{})
	m.models = append(m.models, &model.Quote{})
}

// check table with models, create if not exist
func (m *Manager) CheckTable() {
	m.once.Do(func() {
		for _, model := range m.models {
			if m.db.HasTable(model) {
				if err := m.db.AutoMigrate(model).Error; err != nil {
					log.Printf("Auto Migrate Table %s To DB Error: %s", model.TableName(), err.Error())
				}
			} else {
				if err := m.db.CreateTable(model).Error; err != nil {
					log.Printf("Create Table %s To DB Error: %s", model.TableName(), err)
				}
			}
		}
	})
}