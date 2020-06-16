package businessModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// Goods model 商品
type Goods struct {
	model.Model
	GoodsName        string  `json:"goods_name"`            // 商品名称
	GoodsPhoto       string  `json:"goods_photo"`           // 商品图片
	GoodsDescription string  `json:"goods_description"`     // 商品描述
	GoodsListing     bool    `json:"goods_listing"`         // 是否上架
	GoodsPrice       float64 `json:"goods_price"`           // 商品价格
	GoodsUnit        string  `json:"goods_unit"`            // 商品单位 份、杯
	GoodsSort        int64   `json:"goods_sort"`            // 商品排序
	GoodsTypeID      int64   `json:"goods_type_id"`         // 商品种类ID
	GoodsTypeName    string  `json:"goods_type_name"`       // 商品种类名称
	AdminID          int64   `gorm:"INDEX" json:"admin_id"` // 商家管理员ID
}

// 添加商品
func (g *Goods) AddGoods() error {
	db := mysql.GetMysqlDB()
	return db.Create(&g).Error
}

// 修改商品by id
func (g *Goods) UpdateGoodsByID(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&g).Updates(args).Error
}

// 查询商品by goods_type_id
func (g *Goods) QueryGoodsByGoodsTypeID() (goods []Goods) {
	db := mysql.GetMysqlDB()
	db.Where("goods_type_id = ?", g.GoodsTypeID).Find(&goods)
	return
}

// 查询商品by id
func (g *Goods) QueryGoodsByID() error {
	db := mysql.GetMysqlDB()
	return db.First(&g).Error
}

// 删除商品
// 	param ids
//  return error
func (g *Goods) DeleteGoodsByIds(ids []int64) error {
	db := mysql.GetMysqlDB()
	tx := db.Begin()
	if err := tx.Where("admin_id = ? AND id IN (?)", g.AdminID, ids).Delete(&Goods{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 检查商家是否已创建相同的商品名称
//  param goods_name
//  param admin_id
func (g *Goods) QueryGoodsExistNameByAdminId() error {
	db := mysql.GetMysqlDB()
	return db.Where("goods_name = ? AND admin_id = ?", g.GoodsName, g.AdminID).First(&g).Error
}

// 批量查询商品
func (g *Goods) QueryGoodsByAdminID(pageSize int, page int) (goods []Goods) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", g.AdminID).Limit(pageSize).Offset((page - 1) * pageSize).Order("goods_sort desc").Find(&goods)
	return
}

// 商品总记录数
func (g *Goods) QueryGoodsCountByAdminID() (count int) {
	db := mysql.GetMysqlDB()
	db.Where("admin_id = ?", g.AdminID).Model(&Goods{}).Count(&count)
	return
}
