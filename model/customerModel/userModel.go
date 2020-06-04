package customerModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
	"github.com/google/uuid"
)

// User model 用户表
type User struct {
	model.Model
	Tel      string    `gorm:"not null;unique" json:"tel"` // 手机号
	Password string    `json:"-"`                          // 密码
	IP       string    `json:"ip"`                         // 登录IP
	Token    string    `json:"token"`                      // 授权令牌
	Avatar   string    `json:"avatar"`                     // 头像
	UUID     uuid.UUID `json:"uuid"`                       // uuid
	Wechat   WeChat    `json:"wechat"`                     // 微信信息
	WechatID int64     `json:"-"`                          // 微信信息ID
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
func (u *User) RegisterUser() error {
	db := mysql.GetMysqlDB()
	return db.Create(&u).Error
}

// 查询用户 by ID
// 	param id
//  return User,error
func (u *User) QueryUserByID() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ? ", u.ID).First(&u).Error
}

// 查询用户 by tel
// 	param tel
//  return User,error
func (u *User) QueryUserByTel() error {
	db := mysql.GetMysqlDB()
	return db.Where("tel = ? ", u.Tel).First(&u).Error
}

// 查询用户 by token
// 	param token
//  return User,error
func (u *User) QueryUserByToken() error {
	db := mysql.GetMysqlDB()
	return db.Where("token = ? AND ISNULL(token)=0 AND LENGTH(trim(token))>0", u.Token).First(&u).Error
}

// 更新用户信息
func (u *User) UpdateUser(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&u).Updates(args).Error
}
