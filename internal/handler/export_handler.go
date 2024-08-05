package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ExpertHandlers(e *echo.Echo, db *gorm.DB) {

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
		db.Create(&user)

		return c.JSON(http.StatusOK, u)
	})

	e.GET("/expert", func(c echo.Context) error {
		data := []model.Expert{}
		result := db.Find(&data)

		log.Println(result, data)
		return c.JSON(http.StatusOK, data)
	})

	e.GET("/expert/text", func(c echo.Context) error {
		users := []model.Expert{}
		result := db.Find(&users)

		log.Println(result, users)
		uList := []string{}
		for _, row := range users {
			uList = append(uList, row.Name)
		}
		us := strings.Join(uList[:], ",")
		return c.String(http.StatusOK, us)
	})

}
