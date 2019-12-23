package app

import "github.com/jasonradcliffe/freshness-countdown-api/controllers/ping"

func mapRoutes() {
	router.GET("/ping", ping.Ping)
}
