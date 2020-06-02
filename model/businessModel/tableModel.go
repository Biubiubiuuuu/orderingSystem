package businessModel

import (
	"errors"

	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// TableSeat model 餐桌
type TableSeat struct {
	model.Model
	Name            string `json:"name"`                  // 餐桌名称
	Sort            int64  `json:"sort"`                  // 餐桌排序
	QRCode          string `json:"QR_code"`               // 餐桌二维码
	Opening         bool   `json:"opening"`               // 是否开台
	DisplayOrNot    bool   `json:"display_or_not"`        // 是否显示
	TableSeatTypeID int64  `json:"table_seat_type_id"`    // 餐桌分类ID
	AdminID         int64  `gorm:"INDEX" json:"admin_id"` // 商家管理员ID
}

// 添加餐桌
func (t *TableSeat) AddTableSeat() error {
	db := mysql.GetMysqlDB()
	return db.Create(&t).Error
}

// 修改餐桌
func (t *TableSeat) UpdateTableSeat(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&t).Updates(args).Error
}

// 查询餐桌
func (t *TableSeat) QueryTableSeat() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ?", t.ID).First(&t).Error
}

// 删除餐桌(可批量)
// 	param id
//  return error
func (t *TableSeat) DeleteTableSeat(ids []int64) error {
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

// 检查商家是否已创建相同的餐桌名称
func (t *TableSeat) QueryTableSeatExistName() error {
	db := mysql.GetMysqlDB()
	return db.Where("name = ? AND admin_id = ?", t.Name, t.AdminID).Error
}

// 批量查询餐桌
func (t *TableSeat) QueryTableSeatsByAdminID(pageSize int, page int) (tableSeats []TableSeat) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Limit(pageSize).Offset((page - 1) * pageSize).Order("sort desc").Find(&tableSeats)
	return
}

// 餐桌总记录数
func (t *TableSeat) QueryTableSeatCountByAdminID() (count int) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", t.AdminID).Model(&TableSeat{}).Count(&count)
	return
}
