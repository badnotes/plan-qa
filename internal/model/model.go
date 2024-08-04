package model

import (
	"time"

	"gorm.io/gorm"
)

// test
type Expert struct {
	gorm.Model
	Code  string `json:"code"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// 店铺
type Shop struct {
	gorm.Model
	Boss_phone string `json:"boss_phone"`
	Name       string `json:"name"`
	Info       string `json:"info"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

// 资源
type Resource struct {
	gorm.Model
	Shop_id uint   `json:"shop_id"`
	Name    string `json:"name"`
	Info    string `json:"info"`
	Phone   string `json:"phone"`
}

// 排班
type Scheduling struct {
	gorm.Model
	Shop_id     uint      `json:"shop_id"`
	Sc_date     string    `json:"sc_date"`
	Time_start  time.Time `json:"time_start"`
	Time_end    time.Time `json:"time_end"`
	Time_long   uint      `json:"time_long"`
	Resource_id uint      `json:"resource_id"`
	Occupied    uint      `json:"occupied"` // 是否被预定
}

// 预约
type Appointment struct {
	gorm.Model
	Shop_id      uint   `json:"shop_id"`
	Sc_id        uint   `json:"sc_id"`
	Ap_type      uint   `json:"ap_type"`
	Client_wx    string `json:"client_wx"`
	Client_name  string `json:"client_name"`
	Client_phone string `json:"client_phone"`
}
