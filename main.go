package main

import (
	"github.com/Biubiubiuuuu/orderingSystem/db/mysql"
	"github.com/Biubiubiuuuu/orderingSystem/db/redis"
	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/Biubiubiuuuu/orderingSystem/helper/encryptHelper"
	"github.com/Biubiubiuuuu/orderingSystem/model/businessModel"
	"github.com/Biubiubiuuuu/orderingSystem/model/systemModel"
	"github.com/Biubiubiuuuu/orderingSystem/router"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	//初始化mysql
	mysql.DB.InitConn()
	db := mysql.GetMysqlDB()
	//自动迁移模型
	db.AutoMigrate(&systemModel.SystemAdmin{}, &businessModel.BusinessAdmin{})
	// 添加默认管理员 username:Admin,password:123456
	a := systemModel.SystemAdmin{Username: "admin", Password: encryptHelper.EncryptMD5To32Bit("123456"), Manager: "Y"}
	if err := a.QuerySystemAdminByUsername(); err != nil {
		a.AddSystemAdmin()
	}
	//初始化redis
	redis.DB.InitConn()
	//初始化路由
	router := router.Init()
	//启动
	router.Run(configHelper.HTTPPort)
}
