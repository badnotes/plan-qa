package handler

import "github.com/labstack/echo/v4"

func Parse_shop(c echo.Context) (string, error) {
	sk := c.Request().Header.Get("sk")
	return sk, nil
}
