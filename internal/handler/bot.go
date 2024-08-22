package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/badnotes/plan-qa/internal/model"
	"github.com/labstack/echo/v4"
)

var (
	counter int
	mutex   sync.Mutex
)

func BotHandlers(e *echo.Echo) {
	e.GET("/bot/resource", botGet)
	e.GET("/bot/scheduling", botListScheduling)
	e.POST("/bot/scheduling/t", botScheduling)
	e.POST("/bot/scheduling", botSchedulingByTemplate)
}

func botGet(c echo.Context) error {
	sk, _ := Parse_shop(c)
	data := []model.Resource{}
	result := model.MyDB.Where("sk = ?", sk).Find(&data)
	log.Println(result, data)

	dl := []ResourceDto{}
	for _, row := range data {
		dl = append(dl, ResourceDto{Name: row.Name, Info: row.Info})
	}

	return c.JSON(http.StatusOK, dl)
}

func botListScheduling(c echo.Context) error {
	sk, _ := Parse_shop(c)
	date := c.QueryParams().Get("date")
	log.Println("date:", date)
	if date == "" {
		if time.Now().Hour() >= 18 {
			// 计算明天的时间
			tomorrow := time.Now().AddDate(0, 0, 1)
			// 格式化明天的日期为字符串
			date = tomorrow.Format("2006-01-02")
		} else {
			// 格式化日期为字符串
			date = time.Now().Format("2006-01-02")
		}
	}

	data := []model.Scheduling{}
	result := model.MyDB.Distinct("sc_date", "time_start", "time_long").Where("sk = ? and sc_date = ? and occupied = 0", sk, date).Find(&data)
	log.Println(result, data)

	// res := []model.Resource{}
	// rs := model.MyDB.Where("sk = ?", sk).Find(&res)
	// log.Println(rs, res)
	// resMap := map[uint]string{}
	// for _, row := range res {
	// 	resMap[row.ID] = row.Name
	// }

	dl := []BotScheduling{}
	for _, row := range data {
		dl = append(dl, BotScheduling{
			Sc_date:    row.Sc_date.Format("2006-01-02"),
			Time_start: row.Time_start,
			Time_long:  row.Time_long,
			// Occupied:   row.Occupied,
		})
	}

	return c.JSON(http.StatusOK, dl)
}

func botScheduling(c echo.Context) error {
	sk, _ := Parse_shop(c)
	if sk == "" {
		return c.JSON(http.StatusOK, Status{Code: 0, Msg: "操作失败"})
	}

	u := new(ReqScheduling)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	// default date is tomorrow
	if u.Date == "" {
		if time.Now().Hour() >= 18 {
			// 计算明天的时间
			tomorrow := time.Now().AddDate(0, 0, 1)
			// 格式化明天的日期为字符串
			u.Date = tomorrow.Format("2006-01-02")
		} else {
			// 格式化日期为字符串
			u.Date = time.Now().Format("2006-01-02")
		}
	}
	log.Println("scheduling param: {}, {}", u.Name, u.Time)

	res := []model.Resource{}
	model.MyDB.Where("sk = ?", sk).Find(&res)
	log.Printf("%+v\n", res)
	resMap := map[uint]string{}
	resId := uint(0)
	for _, row := range res {
		if row.Name == u.Name {
			resId = row.ID
		}
		resMap[row.ID] = row.Name
	}
	log.Println("find resource id: ", resId)

	// lock
	mutex.Lock()
	defer mutex.Unlock()

	data := []model.Scheduling{}
	model.MyDB.Where("sk = ? and sc_date = ? and occupied = 0", sk, u.Date).Find(&data)
	log.Printf("%+v\n", data)
	if len(data) == 0 {
		return c.JSON(http.StatusOK, Status{Code: 1, Msg: u.Date + ",已约满"})
	}

	op_sc_id := uint(0)
	op_dt := ""
	for _, row := range data {
		et := StringToUint(strings.Split(row.Time_start, ":")[0])
		log.Println("可预约时间为:", et)
		if u.Time == et && row.Occupied == 0 {
			op_sc_id = row.ID
			op_dt = row.Time_start
			break
		}
	}
	if op_sc_id <= 0 {
		// 使用 map 去重
		uniqueMap := make(map[string]bool)
		for _, row := range data {
			uniqueMap[row.Time_start] = true
		}

		// 将去重后的元素放入一个新的切片
		var uniqueArr []string
		for item := range uniqueMap {
			uniqueArr = append(uniqueArr, item)
		}

		// 将去重后的元素连接成一个字符串
		result := strings.Join(uniqueArr, ", ")

		return c.JSON(http.StatusOK, Status{Code: 1, Msg: "预约失败,目前可预约时间为:" + u.Date + ", " + result})
	}

	model.MyDB.Model(&model.Scheduling{}).Where("id = ?", op_sc_id).Update("occupied", 1)
	counter++
	log.Println("Counter:", counter)

	return c.JSON(http.StatusOK, Status{Code: 0, Msg: "预约成功，时间为：" + u.Date + ", " + op_dt})
}

func botSchedulingByTemplate(c echo.Context) error {
	sk, _ := Parse_shop(c)
	if sk == "" {
		return c.JSON(http.StatusOK, Status{Code: 0, Msg: "操作失败"})
	}

	u := new(ReqBotScheduling)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	log.Println("scheduling param: ", u)

	// default date is tomorrow
	if u.Time == "" && u.Date != "" {
		u.Time = parse_time(u.Date)
	}
	if u.Time == "" && u.Client_time != "" {
		u.Time = parse_time(u.Client_time)
	}
	if u.Date != "" {
		u.Date = parse_date(u.Date)
	}
	if u.Date == "" && u.Client_time != "" {
		u.Date = parse_date(u.Client_time)
	}
	if u.Date == "" {
		// 计算明天的时间
		tomorrow := time.Now().AddDate(0, 0, 1)
		// 格式化明天的日期为字符串
		u.Date = tomorrow.Format("2006-01-02")
	}
	log.Println("scheduling param: ", u.Name, u.Date, u.Time)

	res := []model.Resource{}
	model.MyDB.Where("sk = ?", sk).Find(&res)
	log.Printf("%+v\n", res)
	resMap := map[uint]string{}
	resId := uint(0)
	for _, row := range res {
		if row.Name == u.Name {
			resId = row.ID
		}
		resMap[row.ID] = row.Name
	}
	log.Println("find resource id: ", resId)

	// lock
	mutex.Lock()
	defer mutex.Unlock()

	data := []model.Scheduling{}
	model.MyDB.Where("sk = ? and sc_date = ? and occupied = 0", sk, u.Date).Find(&data)
	log.Printf("%+v\n", data)
	if len(data) == 0 {
		return c.JSON(http.StatusOK, Status{Code: 1, Msg: u.Date + ",已约满"})
	}

	op_sc_id := uint(0)
	op_dt := ""
	for _, row := range data {
		et := StringToUint(strings.Split(row.Time_start, ":")[0])
		ut := StringToUint(strings.Split(u.Time, ":")[0])
		log.Println("可预约时间为:", et)
		if ut == et && row.Occupied == 0 {
			op_sc_id = row.ID
			op_dt = row.Time_start
			break
		}
	}
	if op_sc_id <= 0 {
		// 使用 map 去重
		uniqueMap := make(map[string]bool)
		for _, row := range data {
			uniqueMap[row.Time_start] = true
		}

		// 将去重后的元素放入一个新的切片
		var uniqueArr []string
		for item := range uniqueMap {
			uniqueArr = append(uniqueArr, item)
		}

		// 将去重后的元素连接成一个字符串
		result := strings.Join(uniqueArr, ", ")

		return c.JSON(http.StatusOK, Status{Code: 1, Msg: "预约失败,目前可预约时间为:" + u.Date + ", " + result})
	}

	ap := model.Appointment{
		Sk:                sk,
		Sc_id:             op_sc_id,
		Ap_type:           1,
		Sc_date:           u.Date,
		Time_start:        u.Time,
		Client_wx:         "",
		Client_name:       u.Name,
		Client_phone:      u.Client_phone,
		Client_sex:        u.Client_sex,
		Client_hw:         u.Client_hw,
		Client_experience: u.Client_experience,
		Client_time:       u.Client_time,
		Client_style:      u.Client_style,
		Client_content:    u.Client_content,
		Client_address:    u.Client_address,
	}

	model.MyDB.Model(&model.Appointment{}).Save(&ap)
	model.MyDB.Model(&model.Scheduling{}).Where("id = ?", op_sc_id).Update("occupied", 1)
	counter++
	log.Println("Counter:", counter)

	return c.JSON(http.StatusOK, Status{Code: 0, Msg: "预约成功，时间为：" + u.Date + ", " + op_dt})
}

func parse_date(originalDate string) string {

	if strings.Contains(originalDate, "今天") {
		// 格式化明天的日期为字符串
		return time.Now().Format("2006-01-02")
	}

	if strings.Contains(originalDate, "明天") {
		// 计算明天的时间
		tomorrow := time.Now().AddDate(0, 0, 1)
		// 格式化明天的日期为字符串
		return tomorrow.Format("2006-01-02")
	}

	if strings.Contains(originalDate, "后天") {
		// 计算明天的时间
		tomorrow := time.Now().AddDate(0, 0, 2)
		// 格式化明天的日期为字符串
		return tomorrow.Format("2006-01-02")
	}

	dt := reg_parse_date(originalDate)
	if dt != "" {
		return dt
	}

	if time.Now().Hour() >= 18 {
		// 计算明天的时间
		tomorrow := time.Now().AddDate(0, 0, 1)
		// 格式化明天的日期为字符串
		return tomorrow.Format("2006-01-02")
	} else {
		// 格式化日期为字符串
		return time.Now().Format("2006-01-02")
	}
}

func parse_time(originalDate string) string {
	t := strings.ReplaceAll(originalDate, " ", "")
	if strings.Contains(t, "点") {
		idx := strings.Index(t, "点")
		pt := ""
		if idx > 3 {
			h := t[idx-3 : idx]
			log.Println("user time:", h)
			pt = _pt(h, t)
			if len(pt) == 3 {
				pt = pt[1:3]
			}
			if pt != "" {
				return pt
			}
		}
		if idx > 2 {
			h := t[idx-2 : idx]
			log.Println("user time:", h)
			pt = _pt(h, t)
			if pt != "" {
				return pt
			}
		}
		if idx > 1 {
			h := t[idx-1 : idx]
			log.Println("user time:", h)
			pt = _pt(h, t)
			if pt != "" {
				return pt
			}
		}
	}
	return ""
}

func _pt(h string, originalDate string) string {
	if isNumeric(h) {
		// pass
	} else if h == "一" {
		h = "1"
	} else if h == "二" {
		h = "2"
	} else if h == "三" {
		h = "3"
	} else if h == "四" {
		h = "4"
	} else if h == "五" {
		h = "5"
	} else if h == "六" {
		h = "6"
	} else if h == "七" {
		h = "7"
	}
	if isNumeric(h) {
		if strings.Contains(originalDate, "下午") {
			num, err := strconv.Atoi(h)
			if err != nil {
				fmt.Println("转换失败:", err)
			}
			return strconv.Itoa(12 + num)
		}
		return h
	}
	return ""
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func reg_parse_date(text string) string {

	//text := "今天是2023年4月20日，明天是2023年4月21日。"

	reg_list := []string{
		`\d{4}年\d{1,2}月\d{1,2}日`,
		`\d{4}-\d{1,2}-\d{1,2}`,
		`\d{4}年\d{1,2}月\d{1,2}`,
	}
	for _, row := range reg_list {
		// 定义一个正则表达式来匹配日期
		re := regexp.MustCompile(row)

		// 查找所有匹配的日期
		matches := re.FindAllString(text, -1)

		// 输出匹配的日期
		for _, match := range matches {
			fmt.Println("找到的日期:", match)

			// 解析日期
			date, err := time.Parse("2006年1月2日", match)
			if err != nil {
				fmt.Println("解析日期失败:", err)
				continue
			}
			return date.Format("2006-01-02")
			fmt.Println("解析后的日期:", date)
		}
	}
	return ""
}
