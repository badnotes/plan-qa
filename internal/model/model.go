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

type Account struct {
	gorm.Model
	Phone    string `json:"phone"`
	Pwd      string `json:"pwd"`
	Username string `json:"username"`
	Info     string `json:"info"`
}

// 店铺
type Shop struct {
	gorm.Model
	Sk         string `json:"sk"`
	Boss_phone string `json:"boss_phone"`
	Name       string `json:"name"`
	Info       string `json:"info"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

// 资源
type Resource struct {
	gorm.Model
	Sk    string `json:"sk"`
	Name  string `json:"name"`
	Info  string `json:"info"`
	Phone string `json:"phone"`
}

// 排班
type Scheduling struct {
	gorm.Model
	Sk          string    `json:"sk"`
	Sc_date     time.Time `json:"sc_date"`
	Time_start  string    `json:"time_start"`
	Time_end    string    `json:"time_end"`
	Time_long   uint      `json:"time_long"`
	Resource_id uint      `json:"resource_id"`
	Occupied    uint      `json:"occupied"` // 是否被预定
}

// 预约
type Appointment struct {
	gorm.Model
	Sk                string `json:"sk"`
	Sc_id             uint   `json:"sc_id"`
	Ap_type           uint   `json:"ap_type"`
	Client_wx         string `json:"client_wx"`
	Client_name       string `json:"client_name"`
	Client_phone      string `json:"client_phone"`
	Client_sex        string `json:"client_sex"`
	Client_hw         string `json:"client_hw"`
	Client_experience string `json:"client_experience"`
	Client_time       string `json:"client_time"`
	Client_style      string `json:"client_style"`
	Client_content    string `json:"client_content"`
	Client_address    string `json:"client_address"`
}

// id INTEGER PRIMARY KEY AUTOINCREMENT,
// sk TEXT, -- 所属店铺
// sc_id INTEGER, -- 排班ID
// ap_type INTEGER, -- 预约类型 -- AI/人工
// client_wx text, -- 客户微信
// client_name text, -- 客户名称
// client_phone text, -- 客户电话
// client_sex text, -- 客户性别：
// client_hw text, -- 身高/体重：
// client_experience text, -- 有无拍摄经验：
// client_time text, -- 拍摄档期：明天下午三点
// client_style text, -- 理想拍摄风格：
// client_content text, -- 拍摄内容：
// client_address text, -- 拍摄地址：厦门市思明区前埔西三路288号3楼（前埔宝马楼上）
