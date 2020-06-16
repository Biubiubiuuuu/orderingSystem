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
	// 门店
	apiS := router.Group("/api/v1/business/store")
	apiS.Use(jwtMiddleware.JWT())
	{
		apiS.GET("", businessController.QueryBusinessStoreInfo)
		apiS.PUT("", businessController.UpdateBusinessStoreInfo)
	}
	// 商品种类
	apiGT := router.Group("/api/v1/business/goodstype")
	apiGT.Use(jwtMiddleware.JWT())
	{
		apiGT.POST("", businessController.AddGoodsType)
		apiGT.GET(":id", businessController.QueryGoodsTypeByID)
		apiGT.GET("", businessController.QueryGoodsType)
		apiGT.PUT(":id", businessController.UpdateGoodsType)
		apiGT.DELETE(":id", businessController.DeleteGoodsType)
	}
	// 商品种类
	apiGTS := router.Group("/api/v1/business/goodstypes")
	apiGTS.Use(jwtMiddleware.JWT())
	{
		apiGTS.GET("", businessController.QueryGoodsTypeIDAndName)
	}
	// 商品
	apiG := router.Group("/api/v1/business/goods")
	apiG.Use(jwtMiddleware.JWT())
	{
		apiG.POST("", businessController.AddGoods)
		apiG.GET(":id", businessController.QueryGoodsByID)
		apiG.GET("", businessController.QueryGoods)
		apiG.PUT(":id", businessController.UpdateGoods)
		apiG.PUT(":id/:downorup", businessController.DownOrUpGoods)
		apiG.DELETE(":id", businessController.DeleteGoods)
	}
	// 餐桌种类
	apiTT := router.Group("/api/v1/business/tabletype")
	apiTT.Use(jwtMiddleware.JWT())
	{
		apiTT.POST("", businessController.AddTableType)
		apiTT.GET("", businessController.QueryTableType)
		apiTT.GET(":id", businessController.QueryTableTypeByID)
		apiTT.PUT(":id", businessController.UpdateTableType)
		apiTT.DELETE(":id", businessController.DeleteTableType)
	}
	// 餐桌种类
	apiTTS := router.Group("/api/v1/business/tabletypes")
	apiTTS.Use(jwtMiddleware.JWT())
	{
		apiTTS.GET("", businessController.QueryTableTypeIDAndName)
	}
	// 餐桌
	apiTA := router.Group("/api/v1/business/table")
	apiTA.Use(jwtMiddleware.JWT())
	{
		apiTA.POST("", businessController.AddTable)
		apiTA.PUT(":id", businessController.UpdateTable)
		apiTA.GET(":id", businessController.QueryTableByID)
		apiTA.GET("", businessController.QueryTable)
		apiTA.GET(":id/qrcode", businessController.GetTableqrcode)
		apiTA.DELETE(":id", businessController.DeleteTable)
	}
}
