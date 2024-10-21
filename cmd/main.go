package main

import (
	"fmt"
	"log"

	"github.com/devdouglasferreira/stockscrawler/internal"
	"github.com/devdouglasferreira/stockscrawler/internal/data"
)

func main() {

	db, err := data.OpenDBConnection()
	if err != nil {
		log.Fatalf("Failed to open DB connection: %s", err)
	}
	tickers := data.GetActiveTickers(db)

	for _, ticker := range tickers {

		content, err := internal.FetchURL(ticker.SourceUrl)
		if err != nil {
			log.Fatalf("Failed to fetch URL: %s", err)
		}

		stockPrice, _ := internal.ParseHTML(content)
		stockPrice.Ticker = ticker.Ticker

		data.InsertStockPrice(db, stockPrice)
		fmt.Println(stockPrice.Close)
	}
}
