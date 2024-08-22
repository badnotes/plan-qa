package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/google/uuid"
	"github.com/gookit/cache"
	"github.com/labstack/echo/v4"
)

const Cookie_key = "auth"
const cookie_hour = 3

type LoginBody struct {
	Username string `json:"username"` // 手机号
	Password string `json:"password"`
}

func LoginHandlers(e *echo.Group) {
	e.POST("/login", login)
}

func login(c echo.Context) error {

	u := new(LoginBody)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	data := model.Account{}
	result := model.MyDB.Where("phone = ?", u.Username).First(&data)
	log.Println(result, data)

	if u.Password != "" && data.Pwd == u.Password {
		auth := strings.ReplaceAll(uuid.New().String(), "-", "")
		UpdateCacheAndCookie(c, auth, u.Username)
		return c.JSON(http.StatusOK, Status{Code: 0, Msg: "登陆成功"})
	}
	return c.JSON(http.StatusForbidden, Status{Code: 1, Msg: "手机号或密码错误"})
}

func UpdateCacheAndCookie(c echo.Context, auth string, user string) {
	WriteCookie(c, auth)
	cache.Set(auth, user, 24*time.Hour)
}

func WriteCookie(c echo.Context, v string) bool {
	cookie := new(http.Cookie)
	cookie.Name = Cookie_key
	cookie.Value = v
	cookie.Expires = time.Now().Add(cookie_hour * time.Hour)
	c.SetCookie(cookie)
	return true
}

func ReadCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(Cookie_key)
	if err != nil {
		return "", err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return cookie.Value, nil
}

func VerifyCookie(key string, c echo.Context) (bool, error) {
	u := cache.Get(key)
	UpdateCacheAndCookie(c, key, u.(string))
	return u != "", nil
}
