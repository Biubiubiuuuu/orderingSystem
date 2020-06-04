package businessController

import (
	"net/http"
	"strings"

	"github.com/Biubiubiuuuu/orderingSystem/entity"
	"github.com/Biubiubiuuuu/orderingSystem/service/businessService"
	"github.com/gin-gonic/gin"
)

// @Summary 商家注册
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.BusinessLoginOrRegisterRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/register [POST]
func Register(c *gin.Context) {
	req := entity.BusinessLoginOrRegisterRequest{}
	res := entity.ResponseData{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.Register(req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 商家手机验证码登录
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.BusinessLoginOrRegisterRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/codelogin [POST]
func CodeLogin(c *gin.Context) {
	req := entity.BusinessLoginOrRegisterRequest{}
	res := entity.ResponseData{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.CodeLogin(req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 商家账号密码登录
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.BusinessPassLoginRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/passlogin [POST]
func PassLogin(c *gin.Context) {
	req := entity.BusinessPassLoginRequest{}
	res := entity.ResponseData{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.PassLogin(req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 查询商家门店信息
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/store [GET]
// @Security ApiKeyAuth
func QueryBusinessStoreInfo(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.QueryBusinessStoreInfo(token)
	c.JSON(http.StatusOK, res)
}
