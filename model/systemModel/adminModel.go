package systemModel

import (
	"errors"

	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
	"github.com/google/uuid"
)

// Admin model 系统管理员
type Admin struct {
	model.Model
	Username string    `gorm:"not null;unique;size:10" json:"username"`    // 管理员
	Password string    `gorm:"not null;size:30" json:"-"`                  // 密码
	IP       string    `gorm:"size:30" json:"ip"`                          // 登录IP
	Token    string    `gorm:"size:30" json:"token"`                       // 授权令牌
	Manager  string    `gorm:"not null;default:'N';size:1" json:"manager"` // 操作权限 Y | N
	Avatar   string    `gorm:"size:10" json:"avatar"`                      // 头像
	UUID     uuid.UUID `json:"uuid"`
}

// 系统管理员登录
//  param username
//  param password
//  return Admin,error
func (a *Admin) LoginAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Where("username = ? AND password = ?", a.Username, a.Password).Find(&a).Error
}

// 添加系统管理员账号
func (a *Admin) AddAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Create(&a).Error
}

// 更新管理员 by id
// 	param username
// 	param password
// 	param ip
// 	param token
// 	param manager
// 	param avatar
// 	param uuid
//  return Admin,error
func (a *Admin) UpdateAdmin(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&a).Update(args).Error
}

// 删除系统管理员(可批量)
// 	param id
//  return error
func (a *Admin) DeleteAdmin(ids []int64) error {
	db := mysql.GetMysqlDB()
	tx := db.Begin()
	for _, id := range ids {
		if id == 0 {
			return errors.New("id is not 0")
		}
		a.ID = id
		if err := tx.Delete(&a).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 查询系统管理员
// 	param id
// 	param username
// 	param token
//  return Admin,error
func (a *Admin) QueryAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ? OR username = ? OR (token = ? AND token IS NOT NULL)", a.ID, a.Username, a.Token).First(&a).Error
}

// 批量查询系统管理员
//  param username
//  param manager
//  return Admins
func QueryAdmins(args map[string]interface{}, pageSize int, page int) (admins []Admin) {
	db := mysql.GetMysqlDB()
	db.Where(args).Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&admins)
	return
}

// 系统管理员总记录数
// 	param username
//  param manager
//  return count
func QueryAdminsCount(args map[string]interface{}) (count int) {
	db := mysql.GetMysqlDB()
	db.Where(args).Model(&Admin{}).Count(&count)
	return
}
