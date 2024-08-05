package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ResourceDto struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

type SchedulingDto struct {
	Sc_date    string    `json:"date"`
	Time_start time.Time `json:"time"`
	Time_long  uint      `json:"time_long"`
	Resource   string    `json:"resource"`
	Occupied   uint      `json:"occupied"` // 是否被预定
}

func BotHandlers(e *echo.Echo, db *gorm.DB) {

	e.GET("/bot/resource", func(c echo.Context) error {
		sk, _ := Parse_shop(c)
		data := []model.Resource{}
		result := db.Where("sk = ?", sk).Find(&data)
		log.Println(result, data)

		dl := []ResourceDto{}
		for _, row := range data {
			dl = append(dl, ResourceDto{Name: row.Name, Info: row.Info})
		}

		return c.JSON(http.StatusOK, dl)
	})

	e.GET("/bot/scheduling", func(c echo.Context) error {
		sk, _ := Parse_shop(c)
		data := []model.Scheduling{}
		result := db.Where("sk = ?", sk).Find(&data)
		log.Println(result, data)

		res := []model.Resource{}
		rs := db.Where("sk = ?", sk).Find(&res)
		log.Println(rs, res)
		resMap := map[uint]string{}
		for _, row := range res {
			resMap[row.ID] = row.Name
		}

		dl := []SchedulingDto{}
		for _, row := range data {
			dl = append(dl, SchedulingDto{
				Sc_date:    row.Sc_date.Format("2006-01-02"),
				Time_start: row.Time_start,
				Time_long:  row.Time_long,
				Resource:   resMap[row.Resource_id],
			})
		}

		return c.JSON(http.StatusOK, dl)
	})
}
