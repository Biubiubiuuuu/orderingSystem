package businessService

import (
	"fmt"
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
	data["user"] = b
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
	data["user"] = b
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
		"InStorePhotos":          req.InStorePhotos,
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

// 添加商品种类
func AddGoodsType(token string, req entity.GoodsTypeRequest) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "添加失败，token错误，未找到用户信息"
		return
	}
	if req.Name == "" {
		res.Message = "商品种类名称不能为空"
		return
	}
	g := businessModel.GoodsType{
		Name:         req.Name,
		TypeSort:     req.TypeSort,
		DisplayOrNot: req.DisplayOrNot,
		AdminID:      b.ID,
	}
	if err := g.QueryGoodsTypeExistNameByAdminID(); err == nil {
		res.Message = "添加失败，已存在该商品种类名称"
		return
	}
	if err := g.AddGoodsType(); err != nil {
		res.Message = "添加失败"
		return
	}
	res.Message = "添加成功"
	res.Status = true
	return
}

// 修改商品种类
func UpdateGoodsType(token string, id int64, req entity.GoodsTypeRequest) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "添加失败，token错误，未找到用户信息"
		return
	}
	if req.Name == "" {
		res.Message = "商品种类名称不能为空"
		return
	}
	g := businessModel.GoodsType{
		AdminID: b.ID,
		Name:    req.Name,
	}
	args := map[string]interface{}{
		"name":           req.Name,
		"type_sort":      req.TypeSort,
		"display_or_not": req.DisplayOrNot,
	}
	if err := g.QueryGoodsTypeExistNameByAdminID(); err != nil {
		if err := g.UpdateGoodsTypeByID(args); err != nil {
			res.Message = "修改失败"
			return
		}
		res.Message = "修改成功"
		res.Status = true
		return
	}
	if id != g.ID {
		res.Message = "修改失败，已存在该商品种类名称"
		return
	}
	if err := g.UpdateGoodsTypeByID(args); err != nil {
		res.Message = "修改失败"
		return
	}
	res.Message = "修改成功"
	res.Status = true
	return
}

// 删除商品种类
func DeleteGoodsType(ids []int64) (res entity.ResponseData) {
	if len(ids) == 0 {
		res.Message = "id 不能为空"
		return
	}
	var count int
	for _, v := range ids {
		g := businessModel.Goods{
			GoodsTypeID: v,
		}
		if arrGoods := g.QueryGoodsByGoodsTypeID(); len(arrGoods) > 0 {
			count += len(arrGoods)
		}
	}
	if count > 0 {
		res.Message = fmt.Sprintf("该商品种类下有%v个商品，无法删除改商品分类", count)
		return
	}
	gt := businessModel.GoodsType{}
	if err := gt.DeleteGoodsTypeByIds(ids); err != nil {
		res.Message = "删除失败"
		return
	}
	res.Status = true
	res.Message = "删除成功"
	return
}

// 查询商品种类By ID
func QueryGoodsTypeByID(id int64) (res entity.ResponseData) {
	g := businessModel.GoodsType{}
	g.ID = id
	if err := g.QueryGoodsTypeByID(); err != nil {
		res.Message = "查询失败，未找到商品种类信息"
		return
	}
	res.Message = "查询成功"
	res.Status = true
	data := make(map[string]interface{})
	data["goodstype"] = g
	res.Data = data
	return
}

// 查询商家商品种类
func QueryGoodsType(token string, pageSize int, page int) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "添加失败，token错误，未找到用户信息"
		return
	}
	g := businessModel.GoodsType{AdminID: b.ID}
	if goodsTypes := g.QueryGoodsTypeByAdminID(pageSize, page); len(goodsTypes) == 0 {
		res.Message = "查询失败，未找到商品种类信息"
		return
	} else {
		res.Message = "查询成功"
		res.Status = true
		data := make(map[string]interface{})
		data["goodstype"] = goodsTypes
		res.Data = data
		return
	}
}

// 查询商家商品种类
func QueryGoodsTypeIDAndNameByAdminID(token string) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "添加失败，token错误，未找到用户信息"
		return
	}
	g := businessModel.GoodsType{AdminID: b.ID}
	if goodsTypes := g.QueryGoodsTypeIDAndNameByAdminID(); len(goodsTypes) == 0 {
		res.Message = "查询失败，未找到商品种类信息"
		return
	} else {
		res.Message = "查询成功"
		res.Status = true
		data := make(map[string]interface{})
		data["goodstype"] = goodsTypes
		res.Data = data
		return
	}
}

// 添加商品
func AddGoods(token string, req entity.GoodsRequest) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "添加失败，token错误，未找到用户信息"
		return
	}
	if req.GoodsName == "" {
		res.Message = "商品名称不能为空"
		return
	}
	gt := businessModel.GoodsType{}
	gt.ID = req.GoodsTypeID
	if err := gt.QueryGoodsTypeByID(); err != nil {
		res.Message = "添加失败，不存在该商品种类"
		return
	}
	g := businessModel.Goods{
		GoodsName:        req.GoodsName,
		GoodsPhoto:       req.GoodsPhoto,
		GoodsDescription: req.GoodsDescription,
		GoodsListing:     req.GoodsListing,
		GoodsPrice:       req.GoodsPrice,
		GoodsUnit:        req.GoodsUnit,
		GoodsSort:        req.GoodsSort,
		GoodsTypeID:      req.GoodsTypeID,
		AdminID:          b.ID,
	}
	if err := g.QueryGoodsExistNameByAdminId(); err == nil {
		res.Message = "添加失败，已存在该商品名称"
		return
	}
	if err := g.AddGoods(); err != nil {
		res.Message = "添加失败"
		return
	}
	res.Message = "添加成功"
	res.Status = true
	return
}

// 修改商品
func UpdateGoods(token string, id int64, req entity.GoodsRequest) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "添加失败，token错误，未找到用户信息"
		return
	}
	if req.GoodsName == "" {
		res.Message = "商品名称不能为空"
		return
	}
	gt := businessModel.GoodsType{}
	gt.ID = req.GoodsTypeID
	if err := gt.QueryGoodsTypeByID(); err != nil {
		res.Message = "添加失败，不存在该商品种类"
		return
	}
	g := businessModel.Goods{
		AdminID:   b.ID,
		GoodsName: req.GoodsName,
	}
	args := map[string]interface{}{
		"goods_name":        req.GoodsName,
		"goods_photo":       req.GoodsPhoto,
		"goods_description": req.GoodsDescription,
		"goods_listing":     req.GoodsListing,
		"goods_price":       req.GoodsPrice,
		"goods_unit":        req.GoodsUnit,
		"goods_sort":        req.GoodsSort,
		"goods_type_id":     req.GoodsTypeID,
	}
	if err := g.QueryGoodsExistNameByAdminId(); err != nil {
		if err := g.UpdateGoodsByID(args); err != nil {
			res.Message = "修改失败"
			return
		}
		res.Message = "修改成功"
		res.Status = true
		return
	}
	if id != g.ID {
		res.Message = "修改失败，已存在该商品种类名称"
		return
	}
	if err := g.UpdateGoodsByID(args); err != nil {
		res.Message = "修改失败"
		return
	}
	res.Message = "修改成功"
	res.Status = true
	return
}

// 删除商品
func DeleteGoods(ids []int64) (res entity.ResponseData) {
	if len(ids) == 0 {
		res.Message = "id 不能为空"
		return
	}
	g := businessModel.Goods{}
	if err := g.DeleteGoodsByIds(ids); err != nil {
		res.Message = "删除失败"
		return
	}
	res.Status = true
	res.Message = "删除成功"
	return
}

// 查询商品 by ID
func QueryGoodsByID(id int64) (res entity.ResponseData) {
	g := businessModel.Goods{}
	g.ID = id
	if err := g.QueryGoodsByID(); err != nil {
		res.Message = "查询失败，未找到商品信息"
		return
	}
	res.Message = "查询成功"
	res.Status = true
	data := make(map[string]interface{})
	data["goods"] = g
	res.Data = data
	return
}

// 分页查询商家商品
func QueryGoods(token string, pageSize int, page int) (res entity.ResponseData) {
	b := businessModel.BusinessAdmin{Token: token}
	if err := b.QueryUserByToken(); err != nil {
		res.Message = "添加失败，token错误，未找到用户信息"
		return
	}
	g := businessModel.Goods{AdminID: b.ID}
	if goods := g.QueryGoodsByAdminID(pageSize, page); len(goods) == 0 {
		res.Message = "查询失败，未找到商品信息"
		return
	} else {
		res.Message = "查询成功"
		res.Status = true
		data := make(map[string]interface{})
		data["goods"] = goods
		res.Data = data
		return
	}
}

// 上架/下架商品
func DownOrUpGoods(id int64, downOrup bool) (res entity.ResponseData) {
	g := businessModel.Goods{}
	g.ID = id
	if err := g.QueryGoodsByID(); err != nil {
		res.Message = "未找到商品信息"
		return
	}
	args := map[string]interface{}{
		"GoodsListing": downOrup,
	}
	message := "上架"
	if !downOrup {
		message = "下架"
	}
	if err := g.UpdateGoodsByID(args); err != nil {
		res.Message = fmt.Sprintf("%v%v", message, "失败，未找到商品信息")
		return
	}
	res.Message = fmt.Sprintf("%v%v", message, "成功")
	res.Status = true
	return
}
