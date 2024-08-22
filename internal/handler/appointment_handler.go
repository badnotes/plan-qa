package handler

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
)

func AppointmentHandlers(e *echo.Group) {
	e.GET("/appintment/list", listAppintment)
}

func listAppintment(c echo.Context) error {
	sk, _ := Parse_shop(c)
	data := []model.Appointment{}
	result := model.MyDB.Where("sk = ? and deleted_at is null", sk).Find(&data)
	log.Println(result, data)

	res := []model.Resource{}
	rs := model.MyDB.Where("sk = ?", sk).Find(&res)
	log.Println(rs, res)
	resMap := map[uint]string{}
	for _, row := range res {
		resMap[row.ID] = row.Name
	}

	return c.JSON(http.StatusOK, &Status{Code: 200, Data: data})
}
