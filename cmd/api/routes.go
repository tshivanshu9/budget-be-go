package main

func (app *Application) routes() {
	apiGroup := app.server.Group("/api")
	publicAuthRoutes := apiGroup.Group("/auth")
	{
		publicAuthRoutes.POST("/register", app.handler.RegisterHandler)
		publicAuthRoutes.POST("/login", app.handler.LoginHandler)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthenticationMiddleware)
	{
		profileRoutes.GET("/authenticated/user", app.handler.GetAuthenticatedUser)
		profileRoutes.PATCH("/update/password", app.handler.ChangeUserPassword)
	}
	app.server.GET("/", app.handler.Healthcheck)
}
