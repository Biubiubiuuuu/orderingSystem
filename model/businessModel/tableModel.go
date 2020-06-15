package businessModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// Table model 餐桌
type Table struct {
	model.Model
	Name          string `json:"name"`                  // 餐桌名称
	Sort          int64  `json:"sort"`                  // 餐桌排序
	QRCode        string `json:"QR_code"`               // 餐桌二维码
	Opening       bool   `json:"opening"`               // 是否开台
	DisplayOrNot  bool   `json:"display_or_not"`        // 是否显示
	TableTypeID   int64  `json:"table_type_id"`         // 餐桌种类ID
	TableTypeName int64  `json:"table_type_name"`       // 餐桌种类名称
	AdminID       int64  `gorm:"INDEX" json:"admin_id"` // 商家管理员ID
}

// 添加餐桌
func (t *Table) AddTable() error {
	db := mysql.GetMysqlDB()
	return db.Create(&t).Error
}

// 修改餐桌
func (t *Table) UpdateTable(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&t).Updates(args).Error
}

// 查询餐桌 By ID
func (t *Table) QueryTableByID() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ?", t.ID).First(&t).Error
}

// 查询餐桌 By TableTypeID
func (t *Table) QueryTableByTableTypeID() (tables []Table) {
	db := mysql.GetMysqlDB()
	db.Where("table_type_id = ?", t.TableTypeID).Find(&tables)
	return
}

// 删除餐桌(可批量)
// 	param id
//  return error
func (t *Table) DeleteTable(ids []int64) error {
	db := mysql.GetMysqlDB()
	tx := db.Begin()
	if err := tx.Where("admin_id = ? AND id IN (?)", t.AdminID, ids).Delete(&Table{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 检查商家是否已创建相同的餐桌名称
func (t *Table) QueryTableExistName() error {
	db := mysql.GetMysqlDB()
	return db.Where("name = ? AND admin_id = ?", t.Name, t.AdminID).First(&t).Error
}

// 批量查询餐桌
func (t *Table) QueryTablesByAdminID(pageSize int, page int) (Tables []Table) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Limit(pageSize).Offset((page - 1) * pageSize).Order("sort desc").Find(&Tables)
	return
}

// 餐桌总记录数
func (t *Table) QueryTableCountByAdminID() (count int) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Model(&Table{}).Count(&count)
	return
}
