
-- 账号
CREATE TABLE "accounts" (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	phone TEXT, -- 所属店铺
	pwd TEXT, -- 归属老板
	username TEXT, -- 店铺名称
	info TEXT, -- 介绍
	created_at DATETIME DEFAULT (datetime('now','localtime')), 
	updated_at DATETIME DEFAULT (datetime('now','localtime')), 
	deleted_at DATETIME);

-- 店铺
CREATE TABLE "shop" (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sk TEXT, -- 所属店铺
	boss_phone TEXT, -- 归属老板
	name TEXT, -- 店铺名称
	info TEXT, -- 介绍
	address text, -- 地址
	phone text, -- 电话
	created_at DATETIME DEFAULT (datetime('now','localtime')), 
	updated_at DATETIME DEFAULT (datetime('now','localtime')), 
	deleted_at DATETIME);

-- 店铺配置
CREATE TABLE "shop_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sk TEXT, -- 所属店铺
	work_start DATETIME, -- 上班时间
	work_end DATETIME, -- 下班时间
	work_shift INTEGER, -- 班次时长-分钟
	created_at DATETIME DEFAULT (datetime('now','localtime')), 
	updated_at DATETIME DEFAULT (datetime('now','localtime')), 
	deleted_at DATETIME);

-- 资源
CREATE TABLE "resources" (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sk TEXT, -- 所属店铺
	name TEXT, -- 资源名称
	info TEXT, -- 介绍
	phone text, -- 电话
	created_at DATETIME DEFAULT (datetime('now','localtime')), 
	updated_at DATETIME DEFAULT (datetime('now','localtime')), 
	deleted_at DATETIME);

-- 排班
CREATE TABLE "schedulings" (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sk TEXT, -- 所属店铺
	sc_date Date, -- 日期
    time_start DATETIME, -- 开始时间
	time_end DATETIME, -- 结束时间
    time_long INTEGER, -- 时长，单位-分钟
    resource_id INTEGER, -- 资源ID
	occupied INTEGER DEFAULT 0, -- 是否被预定
	created_at DATETIME DEFAULT (datetime('now','localtime')), 
	updated_at DATETIME DEFAULT (datetime('now','localtime')), 
	deleted_at DATETIME);

-- 预约
CREATE TABLE "appointment" (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sk TEXT, -- 所属店铺
    sc_id INTEGER, -- 排班ID
    ap_type INTEGER, -- 预约类型 -- AI/人工
    client_wx text, -- 客户微信ID
    client_name text, -- 客户名称 
    client_phone text, -- 客户电话 
	created_at DATETIME DEFAULT (datetime('now','localtime')), 
	updated_at DATETIME DEFAULT (datetime('now','localtime')), 
	deleted_at DATETIME);


