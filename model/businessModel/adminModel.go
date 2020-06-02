package businessModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
	"github.com/google/uuid"
)

// BusinessAdmin model 商家管理员
type BusinessAdmin struct {
	model.Model
	Tel      string    `gorm:"not null;unique;" json:"tel"` // 手机号
	Password string    `json:"-"`                           // 密码
	IP       string    `json:"ip"`                          // 登录IP
	Token    string    `json:"token"`                       // 授权令牌
	Avatar   string    `json:"avatar"`                      // 头像
	UUID     uuid.UUID `json:"uuid"`                        // uuid
	Wechat   WeChat    `json:"wechat"`                      // 微信信息
	WechatID int64     `json:"-"`                           // 微信信息ID
}

type WeChat struct {
	ID        int64
	NickName  string `json:"nick_name"`  // 用户昵称
	AvatarUrl string `json:"avatar_url"` // 用户头像地址
	Gender    int64  `json:"gender"`     // 用户性别
	Province  string `json:"province"`   // 用户省市
	City      string `json:"city"`       // 用户城市
	Country   string `json:"country"`    // 用户国家
}

// 注册
func (a *BusinessAdmin) RegisterBusinessAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Create(&a).Error
}

// 验证手机号是否已注册
//  param tel
//  return error
func (a *BusinessAdmin) VerificationTelRegistered() error {
	db := mysql.GetMysqlDB()
	return db.Where("tel = ?", a.Tel).First(&a).Error
}

// 登录
//  param tel
//  param password
//  return error
func (a *BusinessAdmin) Login() error {
	db := mysql.GetMysqlDB()
	return db.Where("tel = ? AND password = ?", a.Tel, a.Password).First(&a).Error
}

// 查询门店管理员
// 	param id
// 	param username
// 	param token
//  return BusinessAdmin,error
func (a *BusinessAdmin) QueryBusinessAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ? OR tel = ? OR (token = ? AND token IS NOT NULL)", a.ID, a.Tel, a.Token).First(&a).Error
}

// 更新管理员 by id
// 	param tel
// 	param password
// 	param ip
// 	param token
// 	param manager
// 	param avatar
// 	param uuid
//  param stores
//  return BusinessAdmin,error
func (a *BusinessAdmin) UpdateBusinessAdmin(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&a).Updates(args).Error
}
