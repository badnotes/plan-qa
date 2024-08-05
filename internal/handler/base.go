package handler

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Parse_shop(c echo.Context) (string, error) {
	sk := c.Request().Header.Get("sk")
	if sk == "" {
		log.Println("sk is empty:", sk)
	}
	return sk, nil
}

type Status struct {
	Code uint   `json:"code"`
	Msg  string `json:"msg"`
}

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}
