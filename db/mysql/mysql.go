package mysql

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/Biubiubiuuuu/orderingSystem/model/adminModel"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MysqlDatabase struct {
	Mysql *gorm.DB
}

var DB *MysqlDatabase
var once sync.Once

// 初始化mysql连接
func (db *MysqlDatabase) InitCoon() {
	once.Do(func() {
		DB = &MysqlDatabase{
			Mysql: InitMysqlDB(),
		}
	})
}

// 初始化mysql连接池
func InitMysqlDB() *gorm.DB {
	var (
		dbType, dbName, user, password, host, tablePrefix string
	)
	dbType = configHelper.DBMysqlType
	dbName = configHelper.DBMysqlName
	user = configHelper.DBMysqlUser
	password = configHelper.DBMysqlPassword
	host = configHelper.DBMysqlHost
	tablePrefix = configHelper.DBMysqlTablePrefix
	// 数据库连接字符串
	connect := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err := gorm.Open(dbType, connect)
	if err != nil {
		log.Fatal(err)
	}
	// 全局设置表名不可以为复数形式
	db.SingularTable(true)
	// debug模式下，打印sql日志到控制台，方便查询问题
	if configHelper.RunMode == "debug" {
		db.LogMode(true)
	}
	// 表明前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return fmt.Sprintf("%v%v", tablePrefix, defaultTableName)
	}
	SetDBConns(db)
	return db
}

// 设置其他信息
func SetDBConns(db *gorm.DB) {
	MaxOpenConns, _ := strconv.Atoi(configHelper.MaxOpenConns)
	MaxIdleConns, _ := strconv.Atoi(configHelper.MaxIdleConns)
	// 设置最大开放连接数
	db.DB().SetMaxOpenConns(MaxOpenConns)
	// 设置最大空闲连接数
	db.DB().SetMaxIdleConns(MaxIdleConns)
}

// 获取mysql连接池
func GetMysqlDB() *gorm.DB {
	return InitMysqlDB()
}

// 初始化mysql数据库并自动迁移模型
func InitMysql() {
	DB.InitCoon()
	db := GetMysqlDB()
	// 自动迁移模型
	db.AutoMigrate(&adminModel.Admin{})
}
