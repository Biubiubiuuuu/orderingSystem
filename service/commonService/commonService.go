package commonService

import (
	"github.com/Biubiubiuuuu/orderingSystem/entity"
	"github.com/Biubiubiuuuu/orderingSystem/helper/utilsHelper"
	"github.com/Biubiubiuuuu/orderingSystem/model/commonModel"
)

// 生成验证码
func VerificationCode(tel string) (res entity.ResponseData) {
	if tel == "" {
		res.Message = "手机号码不能为空"
		return
	}
	if !utilsHelper.CheckTelFormat(tel) {
		res.Message = "手机号码格式不正确"
		return
	}
	code := utilsHelper.GenValidateCode(6)
	if code == "" {
		res.Message = "获取验证码失败"
		return
	}
	v := commonModel.Verificationcode{
		Tel:        tel,
		CreateTime: utilsHelper.GetTimestamp(),
		Code:       code,
	}
	if err := v.AddVerificationcode(); err != nil {
		res.Message = "获取验证码失败"
		return
	}
	// 短信通知接口
	data := make(map[string]interface{})
	data["code"] = code
	res.Data = data
	res.Message = "获取验证码成功"
	res.Status = true
	return
}
