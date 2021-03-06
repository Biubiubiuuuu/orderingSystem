package businessModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// Store model 门店信息
type Store struct {
	model.Model
	StoreName              string         `gorm:"not null" json:"store_name"`                                          // 门店名称
	StoreAddress           string         `gorm:"not null" json:"store_address"`                                       // 门店详细地址
	StoreLogo              string         `json:"store_logo"`                                                          // 门店logo
	StoreContactName       string         `json:"store_contact_name"`                                                  // 门店联系人姓名
	StoreContactTel        string         `json:"store_contact_tel"`                                                   // 门店联系人电话
	StoreStartBankingHours string         `gorm:"not null" json:"store_start_banking_hours"`                           // 门店开始营业时间
	StoreEndBankingHours   string         `gorm:"not null" json:"store_end_banking_hours"`                             // 门店结束营业时间
	StoreFacePhoto         string         `json:"store_face_photo"`                                                    // 门脸照
	InStorePhotos          []InStorePhoto `gorm:"foreignkey:StoreID;association_foreignkey:ID" json:"in_store_photos"` // 店内照
	AdminID                int64          `gorm:"INDEX" json:"admin_id"`                                               // 商家管理员ID
}

// 店内照
type InStorePhoto struct {
	ID      int64
	Url     string `json:"Url"`                   // 图片保存地址
	StoreID int64  `gorm:"INDEX" json:"store_id"` // 门店ID
}

// 添加门店
func (s *Store) AddStore() error {
	db := mysql.GetMysqlDB()
	return db.Create(&s).Error
}

// 修改门店信息
func (s *Store) UpdateStore(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	db.Model(&s).Association("InStorePhotos").Replace(s.InStorePhotos)
	return db.Model(&s).Update(args).Error
}

// 查询门店信息 by AdminID
func (s *Store) QueryStoreByAdminID() error {
	db := mysql.GetMysqlDB()
	return db.First(&s).Model(&s).Related(&s.InStorePhotos).Where("admin_id = ?", s.AdminID).Find(&s).Error
}

// 查询门店信息 by id
func (s *Store) QueryStoreByID() error {
	db := mysql.GetMysqlDB()
	return db.First(&s).Model(&s).Related(&s.InStorePhotos).Find(&s).Error
}
