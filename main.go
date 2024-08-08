package main

import (
	"net/http"
	"strings"

	"github.com/badnotes/plan-qa/internal/handler"
	"github.com/badnotes/plan-qa/internal/model"
	"github.com/gookit/cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {

	// 注册一个（或多个）缓存驱动
	cache.Register(cache.DvrFile, cache.NewFileCache("data/cache.db"))

	model.InitDB()
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

	// auth
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().RequestURI == "/login" || strings.HasPrefix(c.Request().RequestURI, "/bot")
		},
		KeyLookup: "cookie:" + handler.Cookie_key,
		Validator: handler.VerifyCookie,
	}))
	log.Println(e.AcquireContext().Request())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	handler.LoginHandlers(e)
	handler.BotHandlers(e)
	handler.ExpertHandlers(e)
	handler.ShopHandlers(e)
	handler.ResourceHandlers(e)
	handler.SchedulingHandlers(e)

	if err := e.Start(":1323"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
