package app

func mapRoutes() {
	router.GET("/ping", apiHandler.Ping)
	router.GET("/pong", apiHandler.Pong)

	router.POST("/dishes", apiHandler.GetDishes)
	router.POST("/dishes/dish", apiHandler.HandleDishRequest)
	router.POST("/dishes/dish/:p_id", apiHandler.HandleDishRequest)
	router.POST("/dishes/expired", apiHandler.GetDishesExpired)
	router.POST("/dishes/expiresby/:date", apiHandler.GetDishesExpiresBy)

	router.POST("/storage", apiHandler.HandleStorageRequest)

	router.POST("/users", apiHandler.HandleUsersRequest)

	router.GET("/login", apiHandler.Login)
	router.GET("/oauthlogin", apiHandler.Oauthlogin)
	router.GET("/privacy", Privacy)
	router.GET("/success", apiHandler.LoginSuccess)

}
