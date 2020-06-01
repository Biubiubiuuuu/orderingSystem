package businessModel

import (
	"errors"

	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// GoodsType model 商品分类
type GoodsType struct {
	model.Model
	Name         string `gorm:"not null;unique;size:20" json:"name"` // 分类名称
	TypeSort     int64  `json:"type_sort"`                           // 分类排序
	DisplayOrNot int64  `json:"display_or_not"`                      // 是否显示
}

// 添加商品分类
func (g *GoodsType) AddGoodsType() error {
	db := mysql.GetMysqlDB()
	return db.Create(&g).Error
}

// 修改商品分类
func (g *GoodsType) UpdateGoodsType(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&g).Update(args).Error
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
