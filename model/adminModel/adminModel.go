package adminModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/model"
	"github.com/google/uuid"
)

// Admin model 系统管理员
type Admin struct {
	model.Model
	Username string    `gorm:"not null;unique;size:10" json:"username"`   // 管理员
	Password string    `gorm:"not null;size:255" json:"-"`                // 密码
	IP       string    `gorm:"size:30" json:"ip"`                         // 登录IP
	Token    string    `gorm:"size:255" json:"token"`                     // 授权令牌
	Manage   string    `gorm:"not null;default:'N';size:1" json:"Manage"` // 操作权限 Y | N
	UUID     uuid.UUID `json:"uuid"`
}
