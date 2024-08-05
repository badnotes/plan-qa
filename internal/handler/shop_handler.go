package handler

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
)

func ShopHandlers(e *echo.Echo) {

	e.POST("/shop", func(c echo.Context) (err error) {
		u := new(model.Shop)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		// Load into separate struct for security
		data := model.Shop{
			Boss_phone: u.Boss_phone,
			Name:       u.Name,
			Info:       u.Info,
			Phone:      u.Phone,
		}

		log.Println(data)
		model.MyDB.Create(&data)

		return c.JSON(http.StatusOK, u)
	})

	e.GET("/shop", func(c echo.Context) error {
		data := []model.Shop{}
		result := model.MyDB.Find(&data)

		log.Println(result, data)
		return c.JSON(http.StatusOK, data)
	})

}
