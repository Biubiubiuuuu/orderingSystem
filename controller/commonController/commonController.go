package commonController

import (
	"net/http"

	"github.com/Biubiubiuuuu/orderingSystem/service/commonService"
	"github.com/gin-gonic/gin"
)

// @Summary 获取验证码
// @tags 公共接口
// @Accept  application/json
// @Produce  json
// @Param tel query string true "手机号码"
// @Param createtime query string true "创建时间 格式：yyyy-MM-dd HH:mm:ss"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/verificationcode [GET]
func VerificationCode(c *gin.Context) {
	tel := c.Query("tel")
	createtime := c.Query("createtime")
	res := commonService.VerificationCode(tel,createtime)
	c.Json(http.StatusOK, res)
}
