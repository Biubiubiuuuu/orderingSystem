package businessModel

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/model"
)

// StoreQualificationInfo model 资质信息
type StoreQualificationInfo struct {
	model.Model
	StoreCertificationStatus int64                 `json:"store_certification_status"` // 店铺认证状态  0：已认证 | 1：未认证
	StoreOperationIndustry   string                `json:"shop_operation_industry"`    // 店铺经营行业
	SubjectQualification     SubjectQualification  `json:"subject_qualification"`      // 主体资质
	SubjectQualificationID   int64                 `json:"-"`                          // 主体资质ID
	IndustryQualification    IndustryQualification `json:"industry_qualification"`     // 行业资质
	IndustryQualificationID  int64                 `json:"-"`                          // 行业资质ID
	AdminID                  int64                 `gorm:"INDEX" json:"admin_id"`      // 商家管理员ID
}

// SubjectQualification model 主体资质
type SubjectQualification struct {
	ID               int64
	CertificateType  string `json:"certificate_type"`  // 证书类型
	CertificatePhoto string `json:"certificate_photo"` // 证书照片
}

// IndustryQualification model 行业资质
type IndustryQualification struct {
	ID               int64
	CertificateType  string `json:"certificate_type"`  // 证书类型
	CertificatePhoto string `json:"certificate_photo"` // 证书照片
}

// StoreOpeningPersonalInfo model 开店个人信息
type StoreOpeningPersonalInfo struct {
	ID                int64
	CertificateType   string             `json:"certificate_type"`                                                // 证书类型
	CertificatePhotos []CertificatePhoto `gorm:"foreignkey:StoreOpeningPersonalInfoID;association_foreignkey:ID"` // 身份证照
}

// CertificatePhoto model 身份证照
type CertificatePhoto struct {
	ID                         int64
	Url                        string `json:"Url"`                                         // 图片保存地址
	StoreOpeningPersonalInfoID int64  `gorm:"INDEX" json:"store_opening_personal_info_id"` // 开店个人信息ID
}

// 提交认证
func (s *StoreQualificationInfo) AddStoreCertification() error {
	db := mysql.GetMysqlDB()
	return db.Create(&s).Error
}

// 修改认证信息
func (s *StoreQualificationInfo) UpdateStoreCertification(args map[string]interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&s).Updates(args).Error
}
