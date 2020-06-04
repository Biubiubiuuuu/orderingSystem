package commonModel

import (
	"encoding/json"

	"github.com/Biubiubiuuuu/orderingSystem/db/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// 验证码
type Verificationcode struct {
	Tel        string // 手机号码
	CreateTime int64  // 创建时间戳
	Code       string // 生成的验证码
}

// 添加验证码信息
func (v *Verificationcode) AddVerificationcode() error {
	rs := redis.GetRedisDB()
	jsonData, _ := json.Marshal(&v)
	if _, err := rs.Do("set", v.Tel, jsonData); err != nil {
		return err
	}
	return nil
}

// 获取验证码信息
func (v *Verificationcode) GetVerificationcode() error {
	rs := redis.GetRedisDB()
	if o, err := redigo.Bytes(rs.Do("get", v.Tel)); err == nil {
		json.Unmarshal(o, &v)
		return nil
	} else {
		return err
	}
}

// 删除验证码信息
func (v *Verificationcode) DeleteVerificationcode() {
	rs := redis.GetRedisDB()
	rs.Do("DEL", v.Tel)
}
