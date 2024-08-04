package main

import (
	"log"
	"net/http"

	"github.com/badnotes/plan-qa/internal/handler"
	"github.com/labstack/echo"
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
