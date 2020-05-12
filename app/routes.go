package app

func mapRoutes() {
	router.GET("/ping", apiHandler.Ping)
	router.GET("/pong", apiHandler.Pong)

	router.POST("/dishes", apiHandler.HandleDishesRequest)

	router.POST("/storage", apiHandler.HandleStorageRequest)

	router.GET("/users", apiHandler.GetUsers)
	router.GET("/users/:user_id", apiHandler.GetUserHandler)
	router.POST("/users", apiHandler.CreateUser)
	router.DELETE("/users/:user_id", apiHandler.DeleteUser)

	/*
		router.GET("/storage", apiHandler.GetStorageUnits)
		router.GET("/storage/:storage_id", apiHandler.GetStorageByID)
		router.POST("/storage", apiHandler.CreateStorageUnit)
		router.PATCH("/storage/:storage_id", apiHandler.UpdateStorageUnit)
		router.DELETE("/storage/:storage_id", apiHandler.DeleteStorageUnit)
	*/

	//router.GET("/storage/:storage_id/dishes", apiHandler.GetStorageDishes)

	router.GET("/login", apiHandler.Login)
	router.GET("/oauthlogin", apiHandler.Oauthlogin)
	router.GET("/privacy", Privacy)
	router.GET("/success", apiHandler.LoginSuccess)

}
