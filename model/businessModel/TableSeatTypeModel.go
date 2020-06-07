package businessModel

import (
	"errors"

	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// TableSeatType model 餐桌分类
type TableSeatType struct {
	model.Model
	Name         string `json:"name"`                  // 分类名称
	SeatingMin   int64  `json:"seating_min"`           // 最少可坐人数
	SeatingMax   int64  `json:"seating_max"`           // 最多可坐人数
	DisplayOrNot bool   `json:"display_or_not"`        // 是否显示
	AdminID      int64  `gorm:"INDEX" json:"admin_id"` // 商家管理员ID
}

// 添加餐桌分类
func (t *TableSeatType) AddTableSeatType() error {
	db := mysql.GetMysqlDB()
	return db.Create(&t).Error
}

// 修改餐桌分类
func (t *TableSeatType) UpdateTableSeatType(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&t).Updates(args).Error
}

// 查询餐桌分类
func (t *TableSeatType) QueryTableSeatTypeByID() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ?", t.ID).First(&t).Error
}

// 删除餐桌分类(可批量)
// 	param id
//  return error
func (t *TableSeatType) DeleteTableSeatTypeByIDs(ids []int64) error {
	db := mysql.GetMysqlDB()
	tx := db.Begin()
	for _, id := range ids {
		if id == 0 {
			return errors.New("id is not 0")
		}
		t.ID = id
		if err := tx.Delete(&t).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 检查商家是否已创建相同的餐桌分类名称
func (t *TableSeatType) QueryTableSeatTypeExistName() error {
	db := mysql.GetMysqlDB()
	return db.Where("name = ? AND admin_id = ?", t.Name, t.AdminID).First(&t).Error
}

// 批量查询餐桌分类
func (t *TableSeatType) QueryTableSeatTypesByAdminID(pageSize int, page int) (tableSeatTypes []TableSeatType) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&tableSeatTypes)
	return
}

// 餐桌分类总记录数
func (t *TableSeatType) QueryTableSeatTypeCountByAdminID() (count int) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Model(&TableSeatType{}).Count(&count)
	return
}
