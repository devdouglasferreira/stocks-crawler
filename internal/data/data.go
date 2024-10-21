package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/devdouglasferreira/stockscrawler/internal/models"
	"github.com/go-sql-driver/mysql"
)

func OpenDBConnection() (*sql.DB, error) {

	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDR"),
		DBName: "StockPrices",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil

}

func InsertStockPrice(db *sql.DB, ticker string, open float64, close float64, high float64, low float64, volume int) {
	command := "INSERT INTO DailyStockPrices (Ticker, Open, Close, High, Low, Volume, Date) VALUES (?, ?, ?, ?, ?, ?, CURDATE())"
	_, err := db.Exec(command, ticker, open, close, high, low, volume)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Inserted %f for %s sucessfully\n", close, ticker)
	}

}

func GetActiveTickers(db *sql.DB) []models.TrackedTicker {
	query := "SELECT SourceUrl, Ticker FROM TrackedTickers WHERE Enabled = TRUE"

	rows, err := db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No ticker active found")
		} else {
			log.Fatal(err)
		}
	}
	defer rows.Close()

	tickers := make([]models.TrackedTicker, 0)

	for rows.Next() {
		var sourceUrl string
		var ticker string
		err := rows.Scan(&sourceUrl, &ticker)
		if err != nil {
			log.Fatal(err)
		}

		tickers = append(tickers, models.TrackedTicker{Ticker: ticker, SourceUrl: sourceUrl})
	}

	return tickers
}
