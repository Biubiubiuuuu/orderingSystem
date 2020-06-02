package businessModel

import (
	"errors"

	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// GoodsType model 商品分类
type GoodsType struct {
	model.Model
	Name         string `json:"name"`                  // 分类名称
	TypeSort     int64  `json:"type_sort"`             // 分类排序
	DisplayOrNot bool   `json:"display_or_not"`        // 是否显示
	AdminID      int64  `gorm:"INDEX" json:"admin_id"` // 商家管理员ID
}

// 添加商品分类
func (g *GoodsType) AddGoodsType() error {
	db := mysql.GetMysqlDB()
	return db.Create(&g).Error
}

// 修改商品分类
func (g *GoodsType) UpdateGoodsType(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&g).Updates(args).Error
}

// 查询商品分类
func (g *GoodsType) QueryGoodsType() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ?", g.ID).First(&g).Error
}

// 删除商品分类(可批量)
// 	param id
//  return error
func (g *GoodsType) DeleteGoodsType(ids []int64) error {
	db := mysql.GetMysqlDB()
	tx := db.Begin()
	for _, id := range ids {
		if id == 0 {
			return errors.New("id is not 0")
		}
		g.ID = id
		if err := tx.Delete(&g).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 检查商家是否已创建相同的商品分类名称
func (g *GoodsType) QueryGoodsTypeExistName() error {
	db := mysql.GetMysqlDB()
	return db.Where("name = ? AND admin_id = ?", g.Name, g.AdminID).Error
}

// 批量查询商品分类
func (g *GoodsType) QueryGoodsTypeByAdminID(pageSize int, page int) (goodsTypes []GoodsType) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", g.AdminID).Limit(pageSize).Offset((page - 1) * pageSize).Order("type_sort desc").Find(&goodsTypes)
	return
}

// 商品分类总记录数
func (g *GoodsType) QueryGoodsTypeCountByAdminID() (count int) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", g.AdminID).Model(&GoodsType{}).Count(&count)
	return
}
