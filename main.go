package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/michelaquino/api_poc_angular_react/context"
)

func main() {
	apiContext := context.GetAPIContext()
	logger := apiContext.GetLogger()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	logger.Info("Started at port 8888!")
	e.Logger.Fatal(e.Start(":8888"))
}
