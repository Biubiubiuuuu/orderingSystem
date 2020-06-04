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
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/common/verificationcode [GET]
func VerificationCode(c *gin.Context) {
	tel := c.Query("tel")
	res := commonService.VerificationCode(tel)
	c.JSON(http.StatusOK, res)
}
