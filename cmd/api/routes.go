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

	categoryRoutes := apiGroup.Group("/categories", app.appMiddleware.AuthenticationMiddleware)
	{
		categoryRoutes.GET("/all", app.handler.ListCategoriesHandler)
		categoryRoutes.POST("/create", app.handler.CreateCategoryHandler)
		categoryRoutes.DELETE("/delete/:id", app.handler.DeleteCategoryHandler)
	}

	budgetRoutes := apiGroup.Group("/budgets", app.appMiddleware.AuthenticationMiddleware)
	{
		budgetRoutes.POST("/create", app.handler.CreateBudgetHandler)
	}
	app.server.GET("/", app.handler.Healthcheck)
}
