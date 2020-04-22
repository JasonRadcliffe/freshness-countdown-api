package app

import (
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/dishes"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/ping"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/users"
)

func mapRoutes() {
	router.GET("/ping", ping.Ping)

	router.GET("/dishes", dishes.GetDishes)

	//Handles the routes for "/dishes/7" and "/dishes/expired"
	router.GET("/dishes/:dish_id", dishes.GetHandler)

	router.POST("/dishes", dishes.CreateDish)
	router.PATCH("/dishes/:dish_id", dishes.UpdateDish)
	router.DELETE("/dishes/:dish_id", dishes.DeleteDish)

	router.GET("/users", users.GetUsers)
	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
	router.PATCH("/users/:user_id", users.UpdateUser)
	router.DELETE("/users/:user_id", users.DeleteUser)

	router.GET("/storage", storage.GetStorageUnits)
	router.GET("/storage/:storage_id", storage.GetStorageUnit)
	router.POST("/storage", storage.CreateStorageUnit)
	router.PATCH("/storage/:storage_id", storage.UpdateStorageUnit)
	router.DELETE("/storage/:storage_id", storage.DeleteStorageUnit)

	router.GET("/storage/:storage_id/dishes", dishes.GetStorageDishes)

	router.GET("/login", Login)
	router.GET("/oauthlogin", Oauthlogin)
	router.GET("/privacy", Privacy)

}
