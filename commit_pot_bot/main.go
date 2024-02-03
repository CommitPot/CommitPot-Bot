package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "김연우 스트릭좀 관리해봅시다")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
