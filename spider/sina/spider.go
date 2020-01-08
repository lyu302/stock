package sina

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/lyu302/stock/db"
	"github.com/lyu302/stock/db/model"
)

var (
	StocksUrl = "http://money.finance.sina.com.cn/d/api/openapi_proxy.php"
	QuoteUrl  = "http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData"

	UserAgents = []string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)",
	}

	PageCount = 500
	DayCount  = 1024
	SpiderLimit = 60
)

type Spider struct {
	wg sync.WaitGroup
	stocks []*model.Stock
	quotes []*model.Quote
}

type Stock struct {
	Code    int        `json:"code"`
	Day     string     `json:"day"`
	Count   int        `json:"count"`
	Fields  []string   `json:"fields"`
	Items   [][]interface{} `json:"items"`
}

type Quote struct {
	Day    string  `json:"day"`
	Open   string  `json:"open"`
	High   string  `json:"high"`
	Low    string  `json:"low"`
	Close  string  `json:"close"`
	Volume string  `json:"volume"`
}

func (s *Spider) FetchStocks() error {
	var (
		url  string
		page = 1
		pageCount = PageCount
		count = pageCount
		stocks []*Stock
		start = time.Now()
	)

	for count == pageCount {
		count = 0
		url = fmt.Sprintf("%s?__s=[[\"hq\",\"hs_a\",\"\",0,%d,%d]]", StocksUrl, page, pageCount)

		rsp, err := http.Get(url)
		if err != nil {
			return err
		}

		body, _ := ioutil.ReadAll(rsp.Body)
		if err := json.Unmarshal(body, &stocks); err != nil {
			rsp.Body.Close()
			return  err
		}

		if len(stocks) == 0 || len(stocks[0].Items) == 0 {
			return errors.New("spider response is empty")
		}

		for _, stock := range stocks[0].Items {
			symbolCode := stock[0].(string)
			code := stock[1].(string)
			name := stock[2].(string)
			symbol := strings.ReplaceAll(symbolCode, code, "")

			stockModel := &model.Stock{
				IDModel:  model.IDModel{},
				Symbol: symbol,
				Code:   code,
				Name:   name,
				TimeModel: model.TimeModel{},
			}

			s.stocks = append(s.stocks, stockModel)
			count += 1
		}
		
		rsp.Body.Close()

		page += 1
	}

	log.Printf("=== Stocks Spider From Sina Finish: %d, use %dms ===", len(s.stocks), (time.Now().UnixNano() - start.UnixNano())/int64(time.Millisecond))
	return nil
}

func (s *Spider) FetchAndSaveStocks() error {
	if err := s.FetchStocks(); err != nil {
		return err
	}
	
	if err := s.SaveStocksToDb(); err != nil {
		return err
	}

	return nil
}

func (s *Spider) FetchQuotes(dayCount int) error {
	return nil

	if dayCount == 0 {
		dayCount = DayCount
	}

	stocks, err := db.DefaultManager.StockDao().ListAllStocks()
	if err != nil {
		log.Printf("=== Fetch Quote Error On Stocks Query From DB: %s ===", err.Error())
		return err
	}

	//s.wg.Add(len(stocks))

	for _, stock := range stocks {
		//go s.fetchQuoteByStock(stock, dayCount)
		start := time.Now()

		had_quotes, err := db.DefaultManager.QuoteDao().ListQuotesByStockID(int(stock.ID))
		if err != nil || len(had_quotes) > 0 {
			log.Printf("Stock(%s %s) Quote Spider Had Init", stock.Symbol, stock.Code)
			continue
		}
		s.fetchQuoteByStock(stock, dayCount)

		log.Printf("=== Stock(%s %s) Spider Result Save To DB Finish, use %ds ===", stock.Symbol, stock.Code, (time.Now().UnixNano() - start.UnixNano())/int64(time.Second))
		//time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}

	//s.wg.Wait()

	return nil
}

func (s *Spider) FetchLhb() error {
	return nil
}

func (s *Spider) FetchAndSaveQuotes(dayCount int) error {
	if err := s.FetchQuotes(dayCount); err != nil {
		return err
	}

	if err := s.SaveQuotesToDb(); err != nil {
		return err
	}
	return nil
}

func (s *Spider) SaveStocksToDb() error {
	start := time.Now()
	var newStocks []*model.Stock

	for i, stock := range s.stocks {
		symbol := stock.Symbol
		code := stock.Code

		if err := db.DefaultManager.StockDao().AddModel(stock); err != nil {
			log.Printf("Spider Stock (%s %s) Save To Db Error: %s", symbol, code, err)
			return err
		}

		newStocks = s.stocks[i:len(s.stocks)]
	}

	s.stocks = newStocks

	log.Printf("=== Stocks Spider Result Save To DB Finish: %d, use %dms ===", len(s.stocks), (time.Now().UnixNano() - start.UnixNano())/int64(time.Millisecond))
	return nil
}

func (s *Spider) SaveQuotesToDb() error  {
	start := time.Now()
	var newQuotes []*model.Quote

	for i, quote := range s.quotes {
		dateStr := quote.Date.Format("2006-01-02")
		if err := db.DefaultManager.QuoteDao().AddModel(quote); err != nil {
			log.Printf("Spider Stock Quote (%s %s) Save To Db Error: %s", quote.StockID, dateStr, err)
			return err
		}

		newQuotes = s.quotes[i:len(s.quotes)]
	}

	s.quotes = newQuotes

	log.Printf("=== Stock Quotes Spider Result Save To DB Finish: %d, use %dms ===", len(s.quotes), (time.Now().UnixNano() - start.UnixNano())/int64(time.Millisecond))
	return nil
}

func (s *Spider) fetchQuoteByStock(stock *model.Stock, dayCount int) {
	var (
		urlStr       string
		ma        = 5
	)

	if dayCount == 0 {
		dayCount = 1024
	}

	//defer s.wg.Done()

	symbolCode := fmt.Sprintf("%s%s", stock.Symbol, stock.Code)
	urlStr = fmt.Sprintf("%s?symbol=%s&scale=240&ma=%d&datalen=%d", QuoteUrl, symbolCode, ma, dayCount)

	client := http.Client{}
	request, _ := http.NewRequest("GET", urlStr, nil)
	request.Header.Add("User-Agent", UserAgents[rand.Intn(len(UserAgents))])

	rsp, err := client.Do(request) //http.Get(url)
	if err != nil {
		log.Printf("=== Fetch Stock(%s) Quote Error: %s ===", symbolCode, err.Error())
		return
	}

	var quotes []*Quote
	body, _ := ioutil.ReadAll(rsp.Body)
	// modify body format
	body = bytes.ReplaceAll(body, []byte("{"), []byte("{\""))
	body = bytes.ReplaceAll(body, []byte(":"), []byte("\":"))
	body = bytes.ReplaceAll(body, []byte(","), []byte(",\""))
	body = bytes.ReplaceAll(body, []byte(",\"{"), []byte(",{"))

	log.Printf("[%d] %s", len(s.quotes) / dayCount, string(body))
	if err := json.Unmarshal(body, &quotes); err != nil {
		log.Printf("=== Fetch Quote Body Error: %s ===", err.Error())
		rsp.Body.Close()
		return
	}

	var lastClose float64 = -1

	for _, quote := range quotes {
		open, _ := strconv.ParseFloat(quote.Open, 64)
		high, _ := strconv.ParseFloat(quote.High, 64)
		low, _ := strconv.ParseFloat(quote.Low, 64)
		close, _ := strconv.ParseFloat(quote.Close, 64)
		volume, _ := strconv.ParseInt(quote.Volume, 10, 64)
		Date, _ := time.Parse("2006-01-02", quote.Day)

		if lastClose < 0 {
			lastClose = close
			continue
		}

		changePercent := (close - lastClose) * 100 / lastClose
		percent := fmt.Sprintf("%.2f", changePercent)
		changePercent, _ = strconv.ParseFloat(percent, 64)

		q := model.Quote{
			IDModel:       model.IDModel{},
			StockID:       int(stock.ID),
			Open:          open,
			High:          high,
			Low:           low,
			Close:         close,
			Volume:        volume,
			Date:          Date,
			ChangePercent: changePercent,
			Stock:         stock,
		}

		//s.quotes = append(s.quotes, &q)

		if err := db.DefaultManager.QuoteDao().AddModel(&q); err != nil {
			log.Printf("=== Stock Quote (%d %s) Save Error: %s ===", stock.ID, Date)
		}

		lastClose = close
	}

	rsp.Body.Close()
}

func (s *Spider) FixZeroPercent()  {
	zeroQuotes, err := db.DefaultManager.QuoteDao().ListQuotesByZeroPercent()
	if err != nil {
		log.Printf("=== Spider Fix Zero Percent Error: %s", err)
		return
	}

	for _, zeroQuote := range zeroQuotes {
		lastQuote, err := db.DefaultManager.QuoteDao().FindQuoteByLastDate(zeroQuote.StockID, zeroQuote.Date)
		if err != nil {
			log.Printf("=== Spider Fix Zero Percent Not Find Last Quote: %s", err)
			continue
		}

		newPercent := (zeroQuote.Close - lastQuote.Close) / lastQuote.Close
		newPercentStr := fmt.Sprintf("%.2f", newPercent)
		newPercent, _ = strconv.ParseFloat(newPercentStr, 64)

		log.Printf("=== Spider Fix Zero Percent Success: %s %s %f %f %f", lastQuote.Date.Format("2006-01-02"), zeroQuote.Stock.Code, lastQuote.Close, zeroQuote.Close, newPercent)

		//zeroQuote.ChangePercent = newPercent
		//db.DefaultManager.QuoteDao().UpdateModel(zeroQuote)
	}

}