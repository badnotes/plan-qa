package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
)

func ExpertHandlers(e *echo.Echo) {

	e.POST("/expert", func(c echo.Context) (err error) {
		u := new(model.Expert)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		// Load into separate struct for security
		user := model.Expert{
			Code:  u.Code,
			Name:  u.Name,
			Email: u.Email,
		}

		log.Println(user)
		model.MyDB.Create(&user)

		return c.JSON(http.StatusOK, u)
	})

	e.GET("/expert", func(c echo.Context) error {
		data := []model.Expert{}
		result := model.MyDB.Find(&data)

		log.Println(result, data)
		return c.JSON(http.StatusOK, data)
	})

	e.GET("/expert/text", func(c echo.Context) error {
		users := []model.Expert{}
		result := model.MyDB.Find(&users)

		log.Println(result, users)
		uList := []string{}
		for _, row := range users {
			uList = append(uList, row.Name)
		}
		us := strings.Join(uList[:], ",")
		return c.String(http.StatusOK, us)
	})

}
