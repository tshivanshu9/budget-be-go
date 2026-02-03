package main

func (app *Application) routes() {
	app.server.GET("/", app.handler.Healthcheck)
	app.server.POST("/register", app.handler.RegisterHandler)
	app.server.POST("/login", app.handler.LoginHandler)
}
