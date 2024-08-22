package handler

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
)

func SchedulingHandlers(e *echo.Group) {
	e.GET("/scheduling/list", listScheduling)
	e.POST("/scheduling/save", scheduling)
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
			Occupied:   row.Occupied,
		})
	}

	return c.JSON(http.StatusOK, dl)
}

func scheduling(c echo.Context) error {
	sk, _ := Parse_shop(c)
	if sk == "" {
		return c.JSON(http.StatusOK, Status{Code: 0, Msg: "操作失败"})
	}

	u := new(model.Scheduling)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	log.Println("scheduling param: {}, {}", u.ID, u.Occupied)
	model.MyDB.Model(&model.Scheduling{}).Where("id = ?", u.ID).Update("occupied", u.Occupied)

	return c.JSON(http.StatusOK, Status{Code: 0, Msg: "操作成功"})
}
