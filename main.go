package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/michelaquino/api_poc_angular_react/context"
	"github.com/michelaquino/api_poc_angular_react/handlers"
)

func main() {
	apiContext := context.GetAPIContext()
	logger := apiContext.GetLogger()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
	}))

	userHandler := handlers.UserHandler{}

	e.GET("/user", userHandler.GetAllUsers)

	logger.Info("Started at port 8888!")
	e.Logger.Fatal(e.Start(":8888"))
}
