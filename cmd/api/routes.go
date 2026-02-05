package main

func (app *Application) routes() {
	apiGroup := app.server.Group("/api")
	publicAuthRoutes := apiGroup.Group("/auth")
	{
		publicAuthRoutes.POST("/register", app.handler.RegisterHandler)
		publicAuthRoutes.POST("/login", app.handler.LoginHandler)
		publicAuthRoutes.POST("forgot/password", app.handler.ForgotPasswordHandler)
		publicAuthRoutes.POST("reset/password", app.handler.ResetPasswordHandler)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthenticationMiddleware)
	{
		profileRoutes.GET("/authenticated/user", app.handler.GetAuthenticatedUser)
		profileRoutes.PATCH("/update/password", app.handler.ChangeUserPassword)
	}
	app.server.GET("/", app.handler.Healthcheck)
}
