package systemService

import (
	"github.com/Biubiubiuuuu/orderingSystem/entity"
	"github.com/Biubiubiuuuu/orderingSystem/helper/encryptHelper"
	"github.com/Biubiubiuuuu/orderingSystem/helper/jwtHelper"
	"github.com/Biubiubiuuuu/orderingSystem/model/systemModel"
	"github.com/google/uuid"
)

// 系统管理员登录
func Login(req entity.SystemAdminLoginRequest, ip string) (res entity.ResponseData) {
	if req.Username == "" || req.Password == "" {
		res.Message = "登录失败，用户名或密码不能空"
		return
	}
	a := systemModel.SystemAdmin{Username: req.Username, Password: encryptHelper.EncryptMD5To32Bit(req.Password)}
	if err := a.LoginSystemAdmin(); err != nil {
		res.Message = "登录失败，用户名或密码错误"
		return
	}
	token, err := jwtHelper.GenerateToken(req.Username, req.Password)
	if err != nil {
		res.Message = "登录失败，token生成错误！"
		return
	}
	// 写入uuid、token、IP，并返回用户信息
	uuid, _ := uuid.NewUUID()
	args := map[string]interface{}{"token": token, "ip": ip, "uuid": uuid}
	if err := a.UpdateSystemAdmin(args); err != nil {
		res.Message = "登录失败，更新登录信息错误"
		return
	}
	data := make(map[string]interface{})
	data["admin"] = a
	res.Status = true
	res.Message = "登录成功"
	res.Data = data
	return
}

// 添加系统管理员
func Add(token string, req entity.SystemAdminAddRequest) (res entity.ResponseData) {
	if req.Username == "" || req.Password == "" {
		res.Message = "添加失败，用户名或密码不能空"
		return
	}
	a := systemModel.SystemAdmin{Token: token}
	if err := a.QuerySystemAdmin(); err != nil {
		res.Message = "修改失败，token错误，未找到系统管理员信息"
		return
	}
	if req.Manager != "Y" {
		req.Manager = "N"
	}
	newAdmin := systemModel.SystemAdmin{
		Username:  req.Username,
		Password:  encryptHelper.EncryptMD5To32Bit(req.Password),
		Manager:   req.Manager,
		Avatar:    req.Avatar,
		CreatedBy: a.Username,
	}
	if err := newAdmin.QuerySystemAdmin(); err == nil {
		res.Message = "添加失败，用户名已存在"
		return
	}
	if err := newAdmin.AddSystemAdmin(); err != nil {
		res.Message = "添加失败"
		return
	}
	res.Status = true
	res.Message = "添加成功"
	return
}

// 修改密码
func UpdatePass(token string, req entity.SystemAdminUpdatePassRequest) (res entity.ResponseData) {
	if req.NewPassword == "" || req.OldPassword == "" {
		res.Message = "修改失败，新/旧密码不能空"
		return
	}
	a := systemModel.SystemAdmin{Token: token}
	if err := a.QuerySystemAdmin(); err != nil {
		res.Message = "修改失败，token错误，未找到系统管理员信息"
		return
	}
	if a.Password != encryptHelper.EncryptMD5To32Bit(req.OldPassword) {
		res.Message = "修改失败，旧密码错误"
		return
	}
	args := map[string]interface{}{"password": encryptHelper.EncryptMD5To32Bit(req.NewPassword)}
	if err := a.UpdateSystemAdmin(args); err != nil {
		res.Message = "修改失败"
		return
	}
	res.Status = true
	res.Message = "修改成功"
	return
}

// 查询管理员账号
func QueryByLimitOffset(args map[string]interface{}, pageSize int, page int) (res entity.ResponseData) {
	data := make(map[string]interface{})
	data["admins"] = systemModel.QuerySystemAdmins(args, pageSize, page)
	data["count"] = systemModel.QuerySystemAdminsCount(args)
	res.Data = data
	res.Status = true
	res.Message = "查询成功"
	return
}

// 启用/禁用管理员
func IsEnableAdmin(token string, is_enable bool) (res entity.ResponseData) {
	a := systemModel.SystemAdmin{Token: token}
	if err := a.QuerySystemAdmin(); err != nil {
		res.Message = "修改失败，token错误，未找到系统管理员信息"
		return
	}
	args := map[string]interface{}{"is_enable": is_enable}
	if err := a.UpdateSystemAdmin(args); err != nil {
		res.Message = "修改失败"
		return
	}
	res.Status = true
	res.Message = "修改成功"
	return
}

// 删除管理员账号
func DeleteAdmin(ids []int64) (res entity.ResponseData) {
	if len(ids) == 0 {
		res.Message = "id 不能为空"
		return
	}
	admin := systemModel.SystemAdmin{}
	if err := admin.DeleteSystemAdmin(ids); err != nil {
		res.Message = "删除失败"
		return
	}
	res.Status = true
	res.Message = "删除成功"
	return
}

// 修改管理员
func UpdateAdmin(token string, args map[string]interface{}) (res entity.ResponseData) {
	a := systemModel.SystemAdmin{Token: token}
	if err := a.QuerySystemAdmin(); err != nil {
		res.Message = "修改失败，token错误，未找到系统管理员信息"
		return
	}
	if err := a.UpdateSystemAdmin(args); err != nil {
		res.Message = "修改失败"
		return
	}
	res.Status = true
	res.Message = "修改成功"
	return
}

// 查询当前管理员信息
func QueryAdminByToken(token string) (res entity.ResponseData) {
	a := systemModel.SystemAdmin{Token: token}
	if err := a.QuerySystemAdmin(); err != nil {
		res.Message = "修改失败，token错误，未找到系统管理员信息"
		return
	}
	data := make(map[string]interface{})
	data["admin"] = a
	res.Data = data
	res.Status = true
	res.Message = "查询成功"
	return
}
