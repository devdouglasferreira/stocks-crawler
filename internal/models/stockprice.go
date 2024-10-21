package models

type StockPrice struct {
	Ticker string
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume int64
}

type TrackedTicker struct {
	Ticker    string
	SourceUrl string
}
