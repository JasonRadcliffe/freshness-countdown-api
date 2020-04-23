package app

import (
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/dish"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/ping"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/storage"
	"github.com/jasonradcliffe/freshness-countdown-api/controllers/user"
)

func mapRoutes() {
	router.GET("/ping", ping.Ping)

	router.GET("/dishes", dish.GetDishes)

	//Handles the routes for "/dish/7" and "/dish/expired"
	router.GET("/dishes/:dish_id", apiHandler.GetDishHandler)

	router.POST("/dishes", dish.CreateDish)
	router.PATCH("/dishes/:dish_id", dish.UpdateDish)
	router.DELETE("/dishes/:dish_id", dish.DeleteDish)

	router.GET("/users", user.GetUsers)
	router.GET("/users/:user_id", user.GetUser)
	router.POST("/users", user.CreateUser)
	router.PATCH("/users/:user_id", user.UpdateUser)
	router.DELETE("/users/:user_id", user.DeleteUser)

	router.GET("/storage", storage.GetStorageUnits)
	router.GET("/storage/:storage_id", storage.GetStorageUnit)
	router.POST("/storage", storage.CreateStorageUnit)
	router.PATCH("/storage/:storage_id", storage.UpdateStorageUnit)
	router.DELETE("/storage/:storage_id", storage.DeleteStorageUnit)

	router.GET("/storage/:storage_id/dish", dish.GetStorageDishes)

	router.GET("/login", Login)
	router.GET("/oauthlogin", Oauthlogin)
	router.GET("/privacy", Privacy)
	router.GET("/success", LoginSuccess)

}
