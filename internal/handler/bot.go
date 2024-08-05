package handler

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
)

func BotHandlers(e *echo.Echo) {
	e.GET("/bot/resource", botGet)
	e.GET("/bot/scheduling", botListScheduling)
	e.POST("/bot/scheduling", botScheduling)
}

func botGet(c echo.Context) error {
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

func botListScheduling(c echo.Context) error {
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
			Occupied:   row.Occupied,
		})
	}

	return c.JSON(http.StatusOK, dl)
}

func botScheduling(c echo.Context) error {
	sk, _ := Parse_shop(c)
	if sk == "" {
		return c.JSON(http.StatusOK, Status{Code: 0, Msg: "操作失败"})
	}

	u := new(ReqScheduling)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	log.Println("scheduling param: {}, {}", u.Name, u.Time)

	res := []model.Resource{}
	model.MyDB.Where("sk = ?", sk).Find(&res)
	log.Printf("%+v\n", res)
	resMap := map[uint]string{}
	resId := uint(0)
	for _, row := range res {
		if row.Name == u.Name {
			resId = row.ID
		}
		resMap[row.ID] = row.Name
	}
	log.Println("find resource id: ", resId)

	data := []model.Scheduling{}
	model.MyDB.Where("sk = ? and resource_id = ?", sk, resId).Find(&data)
	log.Printf("%+v\n", data)

	sc_id := uint(0)
	for _, row := range data {
		et := StringToUint(row.Time_start.Format("15"))
		if u.Time == et && row.Occupied == 0 {
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
