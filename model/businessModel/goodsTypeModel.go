package businessModel

import (
	"errors"

	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/entity"
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
func (g *GoodsType) UpdateGoodsTypeByID(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&g).Updates(args).Error
}

// 查询商品分类
func (g *GoodsType) QueryGoodsTypeByID() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ?", g.ID).First(&g).Error
}

// 删除商品分类(可批量)
// 	param ids
//  return error
func (g *GoodsType) DeleteGoodsTypeByIds(ids []int64) error {
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
func (g *GoodsType) QueryGoodsTypeExistNameByAdminID() error {
	db := mysql.GetMysqlDB()
	return db.Where("name = ? AND admin_id = ?", g.Name, g.AdminID).First(&g).Error
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

// 查询商品分类ID和名称
func (g *GoodsType) QueryGoodsTypeIDAndNameByAdminID() (res []entity.GoodsTypeResponse) {
	db := mysql.GetMysqlDB()
	db.Table("goods_type").Select("id, name").Where("deleted_at IS NULL AND admin_id = ?", g.AdminID).Scan(&res)
	return
}
