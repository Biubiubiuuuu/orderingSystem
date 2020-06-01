package main

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/db/redis"
	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/Biubiubiuuuu/orderingSystem/model/businessModel"
	"github.com/Biubiubiuuuu/orderingSystem/model/systemModel"
	"github.com/Biubiubiuuuu/orderingSystem/router"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	//初始化mysql
	mysql.DB.InitCoon()
	db := mysql.GetMysqlDB()
	//自动迁移模型
	db.AutoMigrate(&systemModel.Admin{}, &businessModel.Admin{})
	//初始化redis
	redis.InitRedis()
	//初始化路由
	router := router.Init()
	//启动
	router.Run(configHelper.HTTPPort)
}
