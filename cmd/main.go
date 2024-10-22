package main

import (
	"log"
	"os"

	"github.com/devdouglasferreira/stockscrawler/internal"
	"github.com/devdouglasferreira/stockscrawler/internal/data"
)

func main() {

	db, err := data.OpenDBConnection()
	if err != nil {
		log.Fatalf("Failed to open DB connection: %s", err)
	}
	tickers := data.GetActiveTickers(db)

	args := os.Args[0]

	for _, ticker := range tickers {

		content, err := internal.FetchURL(ticker.SourceUrl)
		if err != nil {
			log.Fatalf("Failed to fetch URL: %s", err)
		}

		stockPrice, _ := internal.ParseHTML(content)
		stockPrice.Ticker = ticker.Ticker

		if args == "--daily" {
			data.InsertStockPrice(db, stockPrice)
		} else {
			data.InsertIntraDayStockPrice(db, stockPrice)
		}
	}
}
