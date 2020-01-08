package volume_percent

import (
	"github.com/lyu302/stock/arithmetic"
	"github.com/lyu302/stock/db"
	"github.com/lyu302/stock/db/model"
	"log"
	"math"
	"sync"
	"time"

	"github.com/lyu302/stock/forecast"
)

// 交易量小且涨幅大(涨停)
  
type VolumePercent struct {
	Percent   float64
	Start     time.Time
	End       time.Time
	Duration  int
	Before    int
	After     int

	wg        sync.WaitGroup
	Result    map[int]*forecast.Result
}

func NewVolumePercent() *VolumePercent  {
	return &VolumePercent{
		Percent:  9,
		Start:    time.Now().AddDate(0, 0, -5),
		End:      time.Now().AddDate(0, 0, -1),
		Duration: 1,
		Before:   3,
		After:    5,

		Result: make(map[int]*forecast.Result),
	}
}

func (vp * VolumePercent) Forecast()  {
	stockQuotesMap := make(map[int][]*model.Quote)
	quotes, err := db.DefaultManager.QuoteDao().ListQuotesByTime(vp.Start, vp.End)
	if err != nil {
		log.Printf("=== Forecast With VolumePercent Error: %s", err.Error())
		return
	}
	for _, quote := range quotes {
		if _, existed := stockQuotesMap[quote.StockID]; existed {
			stockQuotesMap[quote.StockID] = append(stockQuotesMap[quote.StockID], quote)
		} else {
			stockQuotesMap[quote.StockID] = []*model.Quote{quote}
		}
	}

	vp.wg.Add(len(stockQuotesMap))

	for _, quotes := range stockQuotesMap {
		go vp.quotesForecast(quotes)
	}

	vp.wg.Wait()
}

func (vp *VolumePercent) quotesForecast(quotes []*model.Quote) {
	defer vp.wg.Done()

	if len(quotes) == 0 {
		return
	}

	stock := quotes[0].Stock
	volPers := make([]float64, 0)

	for _, quote := range quotes {
		vol := float64(quote.Volume)
		per := math.Abs(quote.ChangePercent)
		if per < 1  && per > 0{
			per = 1
		}

		volPer := vol /// per
		volPers = append(volPers, volPer)
		//log.Printf("=== Stock(%s %s) Volume Percent Ratio: %.2f - %.2f - %.2f", stock.Code, quote.Date.Format("2006-01-02"), float64(quote.Volume), quote.ChangePercent, volPer)
	}

	var hits [][]*model.Quote
	//hitIndexes := arithmetic.Sigma3WithDuration(volPers, vp.Duration)
	//hitIndexes := arithmetic.QuantileWithDuration(volPers, vp.Duration)
	hitIndexes := arithmetic.MadWithDuration(volPers, vp.Duration)
	for _, hitIndex := range hitIndexes {
		hitQuotes := make([]*model.Quote, 0)
		if quotes[hitIndex[0]].ChangePercent < vp.Percent {
			continue
		}
		for _, index := range hitIndex {
			quote := quotes[index]
			hitQuotes = append(hitQuotes, quote)
		}
		if len(hitQuotes) > 0 {
			log.Printf("=== Stock(%s %s) Volume Percent Sigma3 Hit : %.2f - %.2f - %d", stock.Code, hitQuotes[0].Date.Format("2006-01-02"), float64(hitQuotes[0].Volume), hitQuotes[0].ChangePercent, len(hitQuotes))
		}
		hits = append(hits, hitQuotes)
	}
}

func (vp *VolumePercent) stockForecast(stock *model.Stock)  {
	quotes, err := db.DefaultManager.QuoteDao().ListQuotesByStockIDAndTime(int(stock.ID), vp.Start, vp.End)
	if err != nil {
		log.Printf("=== Stock Forecast With VolumePercent Error: %s", err.Error())
		return
	}
	//log.Printf("=== Stock(%s) Forecast From %d Quotes, Example: %s", stock.Code, len(quotes), quotes[0].Date.Format("2006-01-02"))
	vp.quotesForecast(quotes)
}