package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Expert struct {
	gorm.Model
	Code  string `json:"code"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	log.Println("db: {}", db.Name())
	if err != nil {
		log.Fatalln("db error: {}", err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/expert", func(c echo.Context) (err error) {
		u := new(Expert)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		// Load into separate struct for security
		user := Expert{
			Code:  u.Code,
			Name:  u.Name,
			Email: u.Email}

		log.Println(user)
		db.Create(&user)

		return c.JSON(http.StatusOK, u)
	})
	e.GET("/expert", func(c echo.Context) error {
		users := []Expert{}
		result := db.Find(&users)

		log.Println(result, users)
		return c.JSON(http.StatusOK, users)
	})

	if err := e.Start(":1323"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
