module github.com/devdouglasferreira/stockscrawler

go 1.22.3

require (
	github.com/go-sql-driver/mysql v1.8.1
	golang.org/x/net v0.28.0
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/devdouglasferreira/stockscrawler/internal => ../internal
