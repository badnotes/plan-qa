package handler

import (
	"log"
	"strconv"
	"time"

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

type ResourceDto struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

type ReqScheduling struct {
	Name string `json:"name"`
	Time uint   `json:"time"`
}

type SchedulingDto struct {
	Sc_date    string    `json:"date"`
	Time_start time.Time `json:"time"`
	Time_long  uint      `json:"time_long"`
	Resource   string    `json:"resource"`
	Occupied   uint      `json:"appointment_status"` // 是否被预定
}

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}
