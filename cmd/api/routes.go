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
		budgetRoutes.GET("/all", app.handler.ListBudgetsHandler)
		budgetRoutes.PUT("/:id/update", app.handler.UpdateBudgetHandler)
		budgetRoutes.DELETE("/:id/delete", app.handler.DeleteBudgetHandler)
	}

	walletRoutes := apiGroup.Group("/wallets", app.appMiddleware.AuthenticationMiddleware)
	{
		walletRoutes.POST("/create", app.handler.CreateWalletHandler)
		walletRoutes.GET("/generate-default", app.handler.GenerateDefaultWalletsHandler)
		walletRoutes.GET("/user-list", app.handler.ListUserWalletsHandler)
	}

	transactionRoutes := apiGroup.Group("/transactions", app.appMiddleware.AuthenticationMiddleware)
	{
		transactionRoutes.POST("/create", app.handler.CreateTransactionHandler)
		transactionRoutes.PUT("/:id/reverse", app.handler.ReverseTransactionHandler)
	}
	app.server.GET("/", app.handler.Healthcheck)
}
