package systemModel

import (
	"errors"
	"strconv"
	"time"

	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
	"github.com/google/uuid"
)

// SystemAdmin model 系统管理员
type SystemAdmin struct {
	model.Model
	Username  string    `gorm:"not null;unique;" json:"username"`           // 用户名
	Nikename  string    `json:"nikename"`                                   // 昵称
	Password  string    `gorm:"not null;" json:"-"`                         // 密码
	IP        string    `json:"ip"`                                         // 登录IP
	Token     string    `json:"token"`                                      // 授权令牌
	Manager   string    `gorm:"not null;default:'N';size:1" json:"manager"` // 操作权限 Y | N
	Avatar    string    `json:"avatar"`                                     // 头像
	CreatedBy string    `json:"created_by"`                                 // 创建人
	IsEnable  bool      `json:"is_enable"`                                  // 是否启用 true| false
	UUID      uuid.UUID `json:"uuid"`
}

// 系统管理员登录
//  param username
//  param password
//  return SystemAdmin,error
func (a *SystemAdmin) LoginSystemAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Where("username = ? AND password = ?", a.Username, a.Password).Find(&a).Error
}

// 添加系统管理员账号
func (a *SystemAdmin) AddSystemAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Create(&a).Error
}

// 更新管理员 by id
func (a *SystemAdmin) UpdateSystemAdmin(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&a).Updates(args).Error
}

// 删除系统管理员(可批量)
// 	param id
//  return error
func (a *SystemAdmin) DeleteSystemAdmin(ids []string) error {
	db := mysql.GetMysqlDB()
	tx := db.Begin()
	for _, id := range ids {
		if id == "" {
			return errors.New("id is not 0")
		}
		v, _ := strconv.ParseInt(id, 10, 64)
		a.ID = v
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
//  return SystemAdmin,error
func (a *SystemAdmin) QuerySystemAdmin() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ? OR username = ? OR (token = ? AND token IS NOT NULL)", a.ID, a.Username, a.Token).First(&a).Error
}

// 批量查询系统管理员
//  param username
//  param manager
//  return SystemAdmins
func QuerySystemAdmins(args map[string]interface{}, pageSize int, page int) (SystemAdmins []SystemAdmin) {
	db := mysql.GetMysqlDB()
	if args["username"] != nil && args["username"].(string) != "" {
		db = db.Where("username LIKE ?", "%"+args["username"].(string)+"%")
	}
	if args["created_at_start"] != nil && args["created_at_end"] != nil {
		created_at_start, _ := time.ParseInLocation("2006-01-02", args["created_at_start"].(string), time.Local)
		created_at_end, _ := time.ParseInLocation("2006-01-02", args["created_at_end"].(string), time.Local)
		if args["created_at_start"].(string) != "" && args["created_at_end"].(string) != "" {
			db = db.Where("created_at BETWEEN ? AND ?", created_at_start, created_at_end)
		} else {
			if args["created_at_start"].(string) != "" {
				created_at_start, _ := time.ParseInLocation("2006-01-02", args["created_at_start"].(string), time.Local)
				db = db.Where("created_at > ?", created_at_start)
			} else if args["created_at_end"].(string) != "" {
				created_at_end, _ := time.ParseInLocation("2006-01-02", args["created_at_end"].(string), time.Local)
				db = db.Where("created_at < ?", created_at_end)
			}
		}
	} else {
		if args["created_at_start"] != nil && args["created_at_start"].(string) != "" {
			created_at_start, _ := time.ParseInLocation("2006-01-02", args["created_at_start"].(string), time.Local)
			db = db.Where("created_at > ?", created_at_start)
		} else if args["created_at_end"] != nil && args["created_at_end"].(string) != "" {
			created_at_end, _ := time.ParseInLocation("2006-01-02", args["created_at_end"].(string), time.Local)
			db = db.Where("created_at < ?", created_at_end)
		}
	}
	if args["manager"] != nil && args["manager"].(string) != "" {
		db = db.Where("manager = ?", args["manager"].(string))
	}
	if args["created_by"] != nil && args["created_by"].(string) != "" {
		db = db.Where("created_by = ?", args["created_by"].(string))
	}
	if _, ok := args["is_enable"]; ok {
		if args["is_enable"] != nil {
			db = db.Where("is_enable = ?", args["is_enable"].(bool))
		}
	}
	db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&SystemAdmins)
	return
}

// 系统管理员总记录数
// 	param username
//  param manager
//  return count
func QuerySystemAdminsCount(args map[string]interface{}) (count int) {
	db := mysql.GetMysqlDB()
	if args["username"] != nil && args["username"].(string) != "" {
		db = db.Where("username LIKE ?", "%"+args["username"].(string)+"%")
	}
	if args["created_at_start"] != nil && args["created_at_end"] != nil {
		created_at_start, _ := time.ParseInLocation("2006-01-02", args["created_at_start"].(string), time.Local)
		created_at_end, _ := time.ParseInLocation("2006-01-02", args["created_at_end"].(string), time.Local)
		if args["created_at_start"].(string) != "" && args["created_at_end"].(string) != "" {
			db = db.Where("created_at BETWEEN ? AND ?", created_at_start, created_at_end)
		} else {
			if args["created_at_start"].(string) != "" {
				created_at_start, _ := time.ParseInLocation("2006-01-02", args["created_at_start"].(string), time.Local)
				db = db.Where("created_at > ?", created_at_start)
			} else if args["created_at_end"].(string) != "" {
				created_at_end, _ := time.ParseInLocation("2006-01-02", args["created_at_end"].(string), time.Local)
				db = db.Where("created_at < ?", created_at_end)
			}
		}
	} else {
		if args["created_at_start"] != nil && args["created_at_start"].(string) != "" {
			created_at_start, _ := time.ParseInLocation("2006-01-02", args["created_at_start"].(string), time.Local)
			db = db.Where("created_at > ?", created_at_start)
		} else if args["created_at_end"] != nil && args["created_at_end"].(string) != "" {
			created_at_end, _ := time.ParseInLocation("2006-01-02", args["created_at_end"].(string), time.Local)
			db = db.Where("created_at < ?", created_at_end)
		}
	}
	if args["manager"] != nil && args["manager"].(string) != "" {
		db = db.Where("manager = ?", args["manager"].(string))
	}
	if args["created_by"] != nil && args["created_by"].(string) != "" {
		db = db.Where("created_by = ?", args["created_by"].(string))
	}
	if _, ok := args["is_enable"]; ok {
		if args["is_enable"] != nil {
			db = db.Where("is_enable = ?", args["is_enable"].(bool))
		}
	}
	db.Model(&SystemAdmin{}).Count(&count)
	return
}
