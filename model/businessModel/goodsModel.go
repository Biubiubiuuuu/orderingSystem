package businessModel

import (
	"errors"

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
	AdminID          int64   `gorm:"INDEX" json:"admin_id"` // 商家管理员ID
}

// 添加商品
func (g *Goods) AddGoods() error {
	db := mysql.GetMysqlDB()
	return db.Create(&g).Error
}

// 修改商品
func (g *Goods) UpdateGoods(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&g).Updates(args).Error
}

// 查询商品
func (g *Goods) QueryGoods() error {
	db := mysql.GetMysqlDB()
	return db.Where("id = ?", g.ID).First(&g).Error
}

// 删除商品(可批量)
// 	param id
//  return error
func (g *Goods) DeleteGoods(ids []int64) error {
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

// 检查商家是否已创建相同的商品名称
func (g *Goods) QueryGoodsExistName() error {
	db := mysql.GetMysqlDB()
	return db.Where("goods_name = ? AND admin_id = ?", g.GoodsName, g.AdminID).Error
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
