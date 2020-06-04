package router

import (
	"github.com/Biubiubiuuuu/orderingSystem/controller/businessController"
	"github.com/Biubiubiuuuu/orderingSystem/controller/commonController"
	"github.com/Biubiubiuuuu/orderingSystem/controller/systemController"
	"github.com/Biubiubiuuuu/orderingSystem/docs"
	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/Biubiubiuuuu/orderingSystem/middleware/crossMiddleware"
	"github.com/Biubiubiuuuu/orderingSystem/middleware/errorMiddleware"
	"github.com/Biubiubiuuuu/orderingSystem/middleware/jwtMiddleware"
	"github.com/Biubiubiuuuu/orderingSystem/middleware/loggerMiddleware"
	"github.com/Biubiubiuuuu/orderingSystem/middleware/systemMiddleware"
	"github.com/gin-gonic/gin"
	ginswagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// 初始化路由
func Init() *gin.Engine {
	// swagger接口文档
	docs.SwaggerInfo.Title = "点餐系统接口"
	docs.SwaggerInfo.Description = "一套完整的点餐系统方案"
	docs.SwaggerInfo.Version = configHelper.Version
	// 设置模式
	if configHelper.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	// 记录日志
	router.Use(loggerMiddleware.Logger())
	// 静态资源路径 /static 开头 或者 取自定义配置
	//router.Static(configHelper.Static, "." + configHelper.Static)
	router.Static("/static", "./static")
	//允许跨域请求
	router.Use(crossMiddleware.Cors())
	// 自定义router
	InitSystemAdmin(router)
	InitCommon(router)
	InitBusiness(router)
	//gin swaager
	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerFiles.Handler))
	//404
	router.NoRoute(errorMiddleware.NotFound)
	return router
}

// 系统管理员
func InitSystemAdmin(router *gin.Engine) {
	api := router.Group("/api/v1/system")
	api.POST("login", systemController.Login)
	api.Use(jwtMiddleware.JWT())
	{
		api.GET("admins", systemController.QueryAdmins)
		// 需要管理权限manager为Y才能操作
		api.Use(systemMiddleware.AdminAuth())
		{
			api.DELETE("admins/:ids", systemController.DeleteAdmins)
		}
	}

	apiA := router.Group("/api/v1/system/admin")
	apiA.Use(jwtMiddleware.JWT())
	{
		apiA.GET("", systemController.QueryAdmin)
		apiA.GET(":id", systemController.QueryAdminByID)
		apiA.PUT("password", systemController.UpdatePass)
		apiA.PUT("", systemController.UpdateAdmin)
		apiA.Use(systemMiddleware.AdminAuth())
		{
			apiA.POST("", systemController.AddAdmin)
			apiA.PUT("enable/:id", systemController.IsEnableAdmin)
			apiA.PUT("manager/:id", systemController.IsManagerAdmin)
			apiA.DELETE(":id", systemController.DeleteAdmin)
		}
	}
}

// 公共接口
func InitCommon(router *gin.Engine) {
	api := router.Group("/api/v1/common")
	api.GET("verificationcode", commonController.VerificationCode)
}

// 商家
func InitBusiness(router *gin.Engine) {
	api := router.Group("/api/v1/business")
	api.POST("register", businessController.Register)
	api.POST("codelogin", businessController.CodeLogin)
	api.POST("passlogin", businessController.PassLogin)
	apiS := router.Group("/api/v1/business/store")
	apiS.Use(jwtMiddleware.JWT())
	{
		apiS.GET("", businessController.QueryBusinessStoreInfo)
	}
}
