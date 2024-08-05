package main

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(sqlite.Open("data/gorm.db"), &gorm.Config{})
	log.Println("db: {}", db.Name())
	if err != nil {
		log.Fatalln("db error: {}", err)
	}
	e := echo.New()

	// Middleware
	log := logrus.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	log.Println(e.AcquireContext().Request())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	handler.BotHandlers(e, db)

	handler.ExpertHandlers(e, db)
	handler.ShopHandlers(e, db)
	handler.ResourceHandlers(e, db)

	if err := e.Start(":1323"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
