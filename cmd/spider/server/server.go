package server

import (
	"github.com/lyu302/stock/cmd/spider/option"
	"github.com/lyu302/stock/db"
	"github.com/lyu302/stock/db/config"
	"github.com/lyu302/stock/forecast/volume_percent"
	"github.com/lyu302/stock/spider/sina"
	"log"
	"math/rand"
	"time"
)

func Run(s *option.Spider) error {

	_, err := db.NewDefaultManager(config.Config{
		DbConnectionInfo: s.Config.DBConnectionInfo,
		DbType:           s.Config.DBType,
	})

	if err != nil {
		log.Printf("Create DB Manager Error Befor Spider: %s", err)
		return err
	}

	spider := &sina.Spider{}
	// init run once, then every week
	//spider.FetchAndSaveStocks()
	spider.FetchAndSaveQuotes(0)
	//spider.FixZeroPercent()

	forecast := volume_percent.NewVolumePercent()
	forecast.Forecast()

	//data := []float64{2.78, 1.79, 4.73, 3.81, 2.78, 1.80, 4.81, 2.79, -666, -666, 1.78, 3.32, 10.8, 10.0}
	//exps := arithmetic.Sigma3WithDuration(data, 1)
	//for _, exp := range exps {
	//	for _, e := range exp {
	//		log.Printf("==== %f ====", data[e])
	//	}
	//}

	return nil
}

func init()  {
	rand.Seed(time.Now().UnixNano())
}