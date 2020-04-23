package app

func mapRoutes() {
	router.GET("/ping", apiHandler.Ping)

	router.GET("/dishes", apiHandler.GetDishes)

	//Handles the routes for "/dish/7" and "/dish/expired"
	router.GET("/dishes/:dish_id", apiHandler.GetDishHandler)

	router.POST("/dishes", apiHandler.CreateDish)
	router.PATCH("/dishes/:dish_id", apiHandler.UpdateDish)
	router.DELETE("/dishes/:dish_id", apiHandler.DeleteDish)

	router.GET("/users", apiHandler.GetUsers)
	router.GET("/users/:user_id", apiHandler.GetUserHandler)
	router.POST("/users", apiHandler.CreateUser)
	router.PATCH("/users/:user_id", apiHandler.UpdateUser)
	router.DELETE("/users/:user_id", apiHandler.DeleteUser)

	router.GET("/storage", apiHandler.GetStorageUnits)
	router.GET("/storage/:storage_id", apiHandler.GetStorageByID)
	router.POST("/storage", apiHandler.CreateStorageUnit)
	router.PATCH("/storage/:storage_id", apiHandler.UpdateStorageUnit)
	router.DELETE("/storage/:storage_id", apiHandler.DeleteStorageUnit)

	router.GET("/storage/:storage_id/dishes", apiHandler.GetStorageDishes)

	router.GET("/login", Login)
	router.GET("/oauthlogin", Oauthlogin)
	router.GET("/privacy", Privacy)
	router.GET("/success", LoginSuccess)

}
