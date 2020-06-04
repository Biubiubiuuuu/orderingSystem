package utilsHelper

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const Time_Format = "2006-01-02 15:04:05"
const TimeFormat = "2006-01-02"

// 验证手机号码格式
// 	bool:符合为 true,否则为 false
func CheckTelFormat(tel string) bool {
	str := `^1[3|4|5|7|8][0-9]{9}$`
	rgx := regexp.MustCompile(str)
	return rgx.MatchString(tel)
}

// 生成验证码
//  param width 验证码位数
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// 获取当前时间戳
func GetTimestamp() int64 {
	t := time.Now()
	return t.Unix()
}

// 时间戳转时间yyyy-MM-dd HH:mm:ss
func TimestampToSTime(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format(Time_Format)
}

// 时间戳转time
func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// 时间yyyy-MM-dd转time
func StringToDTime(strTime string) time.Time {
	time, _ := time.ParseInLocation(TimeFormat, strTime, time.Local)
	return time
}

// 时间yyyy-MM-dd HH:mm:ss转time
func StringToSTime(strTime string) time.Time {
	time, _ := time.ParseInLocation(Time_Format, strTime, time.Local)
	return time
}
