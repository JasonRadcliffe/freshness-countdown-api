package app

func mapRoutes() {
	router.GET("/ping", apiHandler.Ping)
	router.GET("/pong", apiHandler.Pong)

	router.POST("/dishes", apiHandler.HandleDishesRequest)
	router.POST("/dishes/:dish_id", apiHandler.HandleDishesRequest)

	router.POST("/storage", apiHandler.HandleStorageRequest)

	router.POST("/users", apiHandler.HandleUsersRequest)

	router.GET("/login", apiHandler.Login)
	router.GET("/oauthlogin", apiHandler.Oauthlogin)
	router.GET("/privacy", Privacy)
	router.GET("/success", apiHandler.LoginSuccess)

}
