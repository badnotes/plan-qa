package handler

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Parse_shop(c echo.Context) (string, error) {
	sk := c.Request().Header.Get("sk")
	if sk == "" {
		log.Println("sk is empty:", sk)
	}
	return sk, nil
}

type Status struct {
	Code uint   `json:"code"`
	Msg  string `json:"msg"`
}

type ResourceDto struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

type ReqScheduling struct {
	Name string `json:"name"`
	Date string `json:"date"`
	Time uint   `json:"time"`
}

type ReqBotScheduling struct {
	Content           string `json:"content"`
	Name              string `json:"name"`
	Date              string `json:"date"`
	Time              string `json:"time"`
	Client_wx         string `json:"client_wx"`         // 客户微信
	Client_name       string `json:"client_name"`       // 客户名称
	Client_phone      string `json:"client_phone"`      // 客户电话
	Client_sex        string `json:"client_sex"`        // 客户性别：
	Client_hw         string `json:"client_hw"`         // 身高/体重：
	Client_experience string `json:"client_experience"` // 有无拍摄经验：
	Client_time       string `json:"client_time"`       // 拍摄档期：明天下午三点
	Client_style      string `json:"client_style"`      // 理想拍摄风格：
	Client_content    string `json:"client_content"`    // 拍摄内容：
	Client_address    string `json:"client_address"`    // 拍摄地址：厦门市思明区前埔西三路288号3楼（前埔宝马楼上）
}

type BotScheduling struct {
	Sc_date    string `json:"date"`
	Time_start string `json:"time"`
	Time_long  uint   `json:"time_long"`
}

type SchedulingDto struct {
	Sc_date    string `json:"date"`
	Time_start string `json:"time"`
	Time_long  uint   `json:"time_long"`
	Resource   string `json:"resource"`
	Occupied   uint   `json:"appointment_status"` // 是否被预定
}

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}
