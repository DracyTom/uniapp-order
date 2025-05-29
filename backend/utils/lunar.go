package utils

import (
	"time"
)

// 农历转换工具
// 这里提供一个简化的农历转换实现，实际项目中可以使用更完整的农历库

var lunarMonths = []string{
	"正月", "二月", "三月", "四月", "五月", "六月",
	"七月", "八月", "九月", "十月", "十一月", "十二月",
}

var lunarDays = []string{
	"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十",
	"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
	"廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十",
}

type LunarDate struct {
	Year  int    `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
	Text  string `json:"text"`
}

// 简化的农历转换函数（实际项目中建议使用专业的农历库）
func ToLunar(date time.Time) LunarDate {
	// 这里是一个简化实现，实际应该使用准确的农历算法
	// 可以集成 github.com/6tail/lunar-go 等专业库
	
	// 简单的近似计算（仅作演示）
	dayOfYear := date.YearDay()
	lunarMonth := (dayOfYear / 30) % 12
	lunarDay := (dayOfYear % 30)
	if lunarDay == 0 {
		lunarDay = 30
	}
	
	monthStr := lunarMonths[lunarMonth]
	dayStr := lunarDays[lunarDay-1]
	
	return LunarDate{
		Year:  date.Year(),
		Month: monthStr,
		Day:   dayStr,
		Text:  monthStr + dayStr,
	}
}