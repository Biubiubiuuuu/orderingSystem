package businessService

import (
	"time"

	"github.com/Biubiubiuuuu/orderingSystem/entity"
	"github.com/Biubiubiuuuu/orderingSystem/helper/encryptHelper"
	"github.com/Biubiubiuuuu/orderingSystem/helper/jwtHelper"
	"github.com/Biubiubiuuuu/orderingSystem/helper/utilsHelper"
	"github.com/Biubiubiuuuu/orderingSystem/model/businessModel"
	"github.com/Biubiubiuuuu/orderingSystem/model/commonModel"
	"github.com/google/uuid"
)

// 商家注册
func Register(req entity.BusinessLoginOrRegisterRequest) (res entity.ResponseData) {
	if req.Tel == "" {
		res.Message = "手机号码不能为空"
		return
	}
	if req.Code == "" {
		res.Message = "验证码不能为空"
		return
	}
	if !utilsHelper.CheckTelFormat(req.Tel) {
		res.Message = "手机号码格式不正确"
		return
	}
	v := commonModel.Verificationcode{Tel: req.Tel}
	if err := v.GetVerificationcode(); err != nil {
		res.Message = "验证码获取失败"
		return
	}
	t1 := utilsHelper.TimestampToTime(v.CreateTime)
	t2 := time.Now()
	sub := t2.Sub(t1)
	if sub.Seconds() > 60 {
		res.Message = "验证码已过期，请重新获取"
		return
	}
	if v.Code != req.Code {
		res.Message = "验证码错误"
		return
	}
	b := businessModel.BusinessAdmin{Tel: req.Tel}
	if err := b.QueryUserByTel(); err == nil {
		res.Message = "手机号码已注册，请登录"
		return
	}
	if err := b.RegisterBusinessAdmin(); err != nil {
		res.Message = "注册失败"
		return
	}
	res.Status = true
	res.Message = "注册成功"
	return
}

// 商家验证码登录
func CodeLogin(req entity.BusinessLoginOrRegisterRequest, ip string) (res entity.ResponseData) {
	if req.Tel == "" {
		res.Message = "手机号码不能为空"
		return
	}
	if req.Code == "" {
		res.Message = "验证码不能为空"
		return
	}
	if !utilsHelper.CheckTelFormat(req.Tel) {
		res.Message = "手机号码格式不正确"
		return
	}
	b := businessModel.BusinessAdmin{Tel: req.Tel}
	if err := b.QueryUserByTel(); err != nil {
		res.Message = "手机号码未注册，请注册后登录"
		return
	}
	v := commonModel.Verificationcode{Tel: req.Tel}
	if err := v.GetVerificationcode(); err != nil {
		res.Message = "验证码获取失败"
		return
	}
	t1 := utilsHelper.TimestampToTime(v.CreateTime)
	t2 := time.Now()
	sub := t2.Sub(t1)
	if sub.Seconds() > 60 {
		res.Message = "验证码已过期，请重新获取"
		return
	}
	if v.Code != req.Code {
		res.Message = "验证码错误"
		return
	}
	token, err := jwtHelper.GenerateToken(req.Tel, req.Code)
	if err != nil {
		res.Message = "登录失败，token生成错误！"
		return
	}
	// 写入uuid、token、IP，并返回用户信息
	uuid, _ := uuid.NewUUID()
	args := map[string]interface{}{"token": token, "ip": ip, "uuid": uuid}
	if err := b.UpdateBusinessAdmin(args); err != nil {
		res.Message = "登录失败，更新登录信息错误"
		return
	}
	res.Status = true
	res.Message = "登录成功"
	data := make(map[string]interface{})
	data["user"] = b.QueryUserByToken()
	res.Data = data
	return
}

// 商家账号密码登录
func PassLogin(req entity.BusinessPassLoginRequest, ip string) (res entity.ResponseData) {
	if req.Tel == "" || req.Password == "" {
		res.Message = "账号或密码不能为空"
		return
	}
	if !utilsHelper.CheckTelFormat(req.Tel) {
		res.Message = "手机号码格式不正确"
		return
	}
	b := businessModel.BusinessAdmin{
		Tel:      req.Tel,
		Password: encryptHelper.EncryptMD5To32Bit(req.Password),
	}
	if err := b.QueryUserByTel(); err != nil {
		res.Message = "手机号码未注册，请注册后登录"
		return
	}
	if b.Password != req.Password {
		res.Message = "账号或密码错误"
		return
	}
	token, err := jwtHelper.GenerateToken(req.Tel, req.Password)
	if err != nil {
		res.Message = "登录失败，token生成错误！"
		return
	}
	// 写入uuid、token、IP，并返回用户信息
	uuid, _ := uuid.NewUUID()
	args := map[string]interface{}{"token": token, "ip": ip, "uuid": uuid}
	if err := b.UpdateBusinessAdmin(args); err != nil {
		res.Message = "登录失败，更新登录信息错误"
		return
	}
	res.Status = true
	res.Message = "登录成功"
	data := make(map[string]interface{})
	data["user"] = b.QueryUserByToken()
	res.Data = data
	return
}

// 查询商家门店信息
func QueryBusinessStoreInfo(token string) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "查询失败，token错误，未找到用户信息"
		return
	}
	s := businessModel.Store{AdminID: b.ID}
	if err := s.QueryStoreByAdminID(); err != nil {
		res.Message = "查询失败，未找到用户门店信息"
		return
	}
	if err := s.QueryStoreByID(); err != nil {
		res.Message = "查询失败，未找到用户门店信息"
		return
	}
	res.Message = "查询成功"
	res.Status = true
	data := make(map[string]interface{})
	data["store"] = s
	res.Data = data
	return
}

// 更新或添加商家门店信息
func UpdateBusinessStoreInfo(token string, req entity.BusinessStoreRequest) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "更新失败，token错误，未找到用户信息"
		return
	}
	s := businessModel.Store{AdminID: b.ID}
	if err := s.QueryStoreByAdminID(); err != nil {
		store := businessModel.Store{
			StoreName:              req.StoreName,
			StoreAddress:           req.StoreAddress,
			StoreLogo:              req.StoreLogo,
			StoreContactName:       req.StoreContactName,
			StoreContactTel:        req.StoreContactTel,
			StoreEndBankingHours:   req.StoreEndBankingHours,
			StoreStartBankingHours: req.StoreStartBankingHours,
			StoreFacePhoto:         req.StoreFacePhoto,
			AdminID:                b.ID,
		}
		if len(req.InStorePhotos) > 0 {
			var photos []businessModel.InStorePhoto
			for _, k := range req.InStorePhotos {
				photo := businessModel.InStorePhoto{
					Url: k.Url,
				}
				photos = append(photos, photo)
			}
			s.InStorePhotos = photos
		}
		if err := store.AddStore(); err != nil {
			res.Message = "更新失败"
			return
		}
		res.Message = "更新成功"
		res.Status = true
		data := make(map[string]interface{})
		data["store"] = store
		res.Data = data
		return
	}
	args := map[string]interface{}{
		"StoreName":              req.StoreName,
		"StoreAddress":           req.StoreAddress,
		"StoreLogo":              req.StoreLogo,
		"StoreContactName":       req.StoreContactName,
		"StoreContactTel":        req.StoreContactTel,
		"StoreStartBankingHours": req.StoreStartBankingHours,
		"StoreEndBankingHours":   req.StoreEndBankingHours,
		"StoreFacePhoto":         req.StoreFacePhoto,
		"AdminID":                b.ID,
	}
	// 存在修改值为 "" 时，先删除该key
	for k, v := range args {
		switch v.(type) {
		case string:
			if v.(string) == "" {
				delete(args, k)
			}
		}
	}
	if len(req.InStorePhotos) > 0 {
		var photos []businessModel.InStorePhoto
		for _, k := range req.InStorePhotos {
			photo := businessModel.InStorePhoto{
				Url: k.Url,
			}
			photos = append(photos, photo)
		}
		args["InStorePhotos"] = photos
	}
	u := businessModel.Store{}
	u.ID = s.ID
	if err := u.UpdateStore(args); err != nil {
		res.Message = "更新失败"
		return
	}
	res.Message = "更新成功"
	res.Status = true
	data := make(map[string]interface{})
	data["store"] = u
	res.Data = data
	return
}
