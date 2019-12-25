package app

import (
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/dishes"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/ping"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/sunits"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/users"
)

func mapRoutes() {
	router.GET("/ping", ping.Ping)

	router.GET("/dishes", dishes.GetDishes)
	router.GET("/dishes/:dish_id", dishes.GetDish)
	router.POST("/dishes", dishes.CreateDish)
	router.PATCH("/dishes/:dish_id", dishes.UpdateDish)
	router.DELETE("/dishes/:dish_id", dishes.DeleteDish)

	router.GET("/users", users.GetUsers)
	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
	router.PATCH("/users/:user_id", users.UpdateUser)
	router.DELETE("/users/:user_id", users.DeleteUser)

	router.GET("/storageunits", sunits.GetSUnits)
	router.GET("/storageunits/:storage_id", sunits.GetSUnit)
	router.POST("/storageunits", sunits.CreateSUnit)
	router.PATCH("/storageunits/:storage_id", sunits.UpdateSUnit)
	router.DELETE("/storageunits/:storage_id", sunits.DeleteSUnit)

	router.GET("/storageunits/:storage_id/dishes", dishes.GetStorageDishes)
	router.GET("/dishes/expired", dishes.GetExpiredDishes)

}
