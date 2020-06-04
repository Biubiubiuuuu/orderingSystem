package systemMiddleware

import (
	"net/http"
	"strings"

	"github.com/Biubiubiuuuu/orderingSystem/common/responseCommon"
	"github.com/Biubiubiuuuu/orderingSystem/model/systemModel"
	"github.com/gin-gonic/gin"
)

// 根据token查询对应账号信息验证当前登录用户或者外部调用 是否有操作curd的权限
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			authToken := c.GetHeader("Authorization")
			if authToken == "" {
				message := "Query not 'token' param OR header Authorization has not Bearer token"
				res := responseCommon.Response(false, nil, message)
				c.AbortWithStatusJSON(http.StatusUnauthorized, res)
				return
			}
			token = strings.TrimSpace(authToken)
		}
		a := systemModel.SystemAdmin{Token: token}
		if err := a.QuerySystemAdminByToken(); err != nil {
			message := "操作失败，token错误，未找到系统管理员信息"
			res := responseCommon.Response(false, nil, message)
			c.AbortWithStatusJSON(http.StatusOK, res)
			return
		}
		if a.Manager != "Y" {
			message := "操作失败，没有权限操作"
			res := responseCommon.Response(false, nil, message)
			c.AbortWithStatusJSON(http.StatusOK, res)
			return
		}
		c.Next()
	}
}
