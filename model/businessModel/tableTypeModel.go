package businessModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/entity"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// TableType model 餐桌种类
type TableType struct {
	model.Model
	Name         string `json:"name"`                  // 种类名称
	SeatingMin   int64  `json:"seating_min"`           // 最少可坐人数
	SeatingMax   int64  `json:"seating_max"`           // 最多可坐人数
	DisplayOrNot bool   `json:"display_or_not"`        // 是否显示
	AdminID      int64  `gorm:"INDEX" json:"admin_id"` // 商家管理员ID
}

// 添加餐桌种类
func (t *TableType) AddTableType() error {
	db := mysql.GetMysqlDB()
	return db.Create(&t).Error
}

// 修改餐桌种类
func (t *TableType) UpdateTableType(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&t).Updates(args).Error
}

// 查询餐桌种类 by id
func (t *TableType) QueryTableTypeByID() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ?", t.ID).First(&t).Error
}

// 删除餐桌种类(可批量)
// 	param id
//  return error
func (t *TableType) DeleteTableTypeByIDs(ids []int64) error {
	db := mysql.GetMysqlDB()
	tx := db.Begin()
	if err := tx.Where("admin_id = ? AND id IN (?)", t.AdminID, ids).Delete(&TableType{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 检查商家是否已创建相同的餐桌种类名称
func (t *TableType) QueryTableTypeExistName() error {
	db := mysql.GetMysqlDB()
	return db.Where("name = ? AND admin_id = ?", t.Name, t.AdminID).First(&t).Error
}

// 批量查询餐桌种类
func (t *TableType) QueryTableTypesByAdminID(pageSize int, page int) (TableTypes []TableType) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&TableTypes)
	return
}

// 餐桌种类总记录数
func (t *TableType) QueryTableTypeCountByAdminID() (count int) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Model(&TableType{}).Count(&count)
	return
}

// 查询餐桌种类ID和名称
func (t *TableType) QueryTableTypeIDAndNameByAdminID() (res []entity.TableTypeResponse) {
	db := mysql.GetMysqlDB()
	db.Table("table_type").Select("id, name").Where("deleted_at IS NULL AND admin_id = ?", t.AdminID).Scan(&res)
	return
}
