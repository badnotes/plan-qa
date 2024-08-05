package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
)

func BotHandlers(e *echo.Echo) {
	e.GET("/bot/resource", get)
	e.GET("/bot/scheduling", listScheduling)
	e.POST("/bot/scheduling", scheduling)
}

type ResourceDto struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

type SchedulingDto struct {
	Sc_date    string    `json:"date"`
	Time_start time.Time `json:"time"`
	Time_long  uint      `json:"time_long"`
	Resource   string    `json:"resource"`
	Occupied   uint      `json:"appointment_status"` // 是否被预定
}

func get(c echo.Context) error {
	sk, _ := Parse_shop(c)
	data := []model.Resource{}
	result := model.MyDB.Where("sk = ?", sk).Find(&data)
	log.Println(result, data)

	dl := []ResourceDto{}
	for _, row := range data {
		dl = append(dl, ResourceDto{Name: row.Name, Info: row.Info})
	}

	return c.JSON(http.StatusOK, dl)
}

func listScheduling(c echo.Context) error {
	sk, _ := Parse_shop(c)
	data := []model.Scheduling{}
	result := model.MyDB.Where("sk = ?", sk).Find(&data)
	log.Println(result, data)

	res := []model.Resource{}
	rs := model.MyDB.Where("sk = ?", sk).Find(&res)
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
}

func scheduling(c echo.Context) error {
	sk, _ := Parse_shop(c)
	if sk == "" {
		return c.JSON(http.StatusOK, Status{Code: 0, Msg: "操作失败"})
	}
	name := c.QueryParams().Get("name")
	t := c.QueryParams().Get("time")
	log.Println("scheduling param: {}, {}", name, t)

	res := []model.Resource{}
	model.MyDB.Where("sk = ?", sk).Find(&res)
	log.Printf("%+v\n", res)
	resMap := map[uint]string{}
	resId := uint(0)
	for _, row := range res {
		if row.Name == name {
			resId = row.ID
		}
		resMap[row.ID] = row.Name
	}
	log.Println("find resource id: {}", resId)

	data := []model.Scheduling{}
	model.MyDB.Where("sk = ? and resource_id = ?", sk, resId).Find(&data)
	log.Printf("%+v\n", data)

	sc_id := uint(0)
	for _, row := range data {
		if t == row.Time_start.Format("15") && row.Occupied == 0 {
			sc_id = row.ID
			break
		}
	}
	if sc_id <= 0 {
		return c.JSON(http.StatusOK, Status{Code: 0, Msg: "操作失败"})
	}
	model.MyDB.Model(&model.Scheduling{}).Where("id = ?", sc_id).Update("occupied", 1)

	return c.JSON(http.StatusOK, Status{Code: 0, Msg: "预定成功"})
}
