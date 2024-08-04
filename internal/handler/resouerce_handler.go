package handler

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func ResourceHandlers(e *echo.Echo, db *gorm.DB) {

	e.POST("/resource", func(c echo.Context) (err error) {
		u := new(model.Resource)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		// Load into separate struct for security
		data := model.Resource{
			Shop_id: u.Shop_id,
			Name:    u.Name,
			Info:    u.Info,
			Phone:   u.Phone,
		}

		log.Println(data)
		db.Create(&data)

		return c.JSON(http.StatusOK, u)
	})

	e.GET("/resource", func(c echo.Context) error {
		data := []model.Resource{}
		result := db.Find(&data)

		log.Println(result, data)
		return c.JSON(http.StatusOK, data)
	})

}
