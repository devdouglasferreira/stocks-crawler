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

func InsertStockPrice(db *sql.DB, stock *models.StockPrice) {
	command := "INSERT INTO DailyStockPrices (Ticker, Open, Close, High, Low, Volume, `Date`) VALUES (?, ?, ?, ?, ?, ?, CURDATE())"
	_, err := db.Exec(command, stock.Ticker, stock.Open, stock.Close, stock.High, stock.Low, stock.Volume)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Inserted %f for %s sucessfully\n", stock.Close, stock.Ticker)
	}

}

func InsertIntraDayStockPrice(db *sql.DB, stock *models.StockPrice) {
	query := "SELECT Value FROM IntraDayStockPrices WHERE Ticker = ? ORDER BY `DateTime` DESC LIMIT 1"

	var lastStockValue float64
	rows := db.QueryRow(query, stock.Ticker)
	rows.Scan(&lastStockValue)

	if stock.Close != lastStockValue {
		command := "INSERT INTO IntraDayStockPrices (Ticker, Open, Value, High, Low, Volume, `DateTime`) VALUES (?, ?, ?, ?, ?, ?, NOW())"
		_, err := db.Exec(command, stock.Ticker, stock.Open, stock.Close, stock.High, stock.Low, stock.Volume)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Inserted %f for %s sucessfully\n", stock.Close, stock.Ticker)
		}
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
