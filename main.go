package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
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

	wx_url, err := url.Parse("http://110.40.156.115:8000")
	if err != nil {
		e.Logger.Printf(err.Error())
	}

	// /wx转发给bot-ai
	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().RequestURI, "/wx")
		},
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{
				URL: wx_url,
			},
		}),
		Rewrite: map[string]string{
			"^/wx/*": "/wx/$1",
		},
		RegexRewrite: map[*regexp.Regexp]string{
			regexp.MustCompile("^/foo/([0-9].*)"):  "/num/$1",
			regexp.MustCompile("^/bar/(.+?)/(.*)"): "/baz/$2/$1",
		},
	}))

	// auth
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			return true
			// return c.Request().RequestURI == "/login" || strings.HasPrefix(c.Request().RequestURI, "/bot")
		},
		KeyLookup: "cookie:" + handler.Cookie_key,
		Validator: handler.VerifyCookie,
	}))
	log.Println(e.AcquireContext().Request())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// bot plugin
	handler.BotHandlers(e)

	// Group level middleware
	g := e.Group("/api")
	handler.LoginHandlers(g)
	handler.ExpertHandlers(g)
	handler.ShopHandlers(g)
	handler.ResourceHandlers(g)
	handler.SchedulingHandlers(g)
	handler.AppointmentHandlers(g)

	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("reoutes:", string(data))

	if err := e.Start(":1323"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
