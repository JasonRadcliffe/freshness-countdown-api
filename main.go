package main

import (
	"github.com/jasonradcliffe/freshness-countdown-api/app"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app.StartApplication()
}
