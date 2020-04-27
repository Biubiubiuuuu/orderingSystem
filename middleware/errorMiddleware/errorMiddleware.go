package errorMiddleware

import (
	"net/http"

	"github.com/Biubiubiuuuu/orderingSystem/common/responseCommon"
	"github.com/gin-gonic/gin"
)

// 404
func NotFound(c *gin.Context) {
	response := responseCommon.Response(false, nil, "404 Not Found")
	c.JSON(http.StatusNotFound, response)
}
