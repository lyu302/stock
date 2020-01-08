package spider

type Spider interface {
	FetchStocks() error
	FetchLhb() error
	FetchQuotes() error
}
