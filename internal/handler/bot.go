package handler

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type ResourceDto struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

func BotHandlers(e *echo.Echo, db *gorm.DB) {

	e.GET("/bot/resource", func(c echo.Context) error {
		data := []model.Resource{}
		result := db.Find(&data)
		log.Println(result, data)

		dl := []ResourceDto{}
		for _, row := range data {
			dl = append(dl, ResourceDto{Name: row.Name, Info: row.Info})
		}

		return c.JSON(http.StatusOK, dl)
	})
}
