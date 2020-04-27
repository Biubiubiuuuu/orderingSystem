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

// 初始化数据库连接
func (db *MysqlDatabase) Init() {
	once.Do(func() {
		DB = &MysqlDatabase{
			Mysql: InitDB(),
		}
	})
}

// 初始化连接对象
func InitDB() *gorm.DB {
	var (
		dbType, dbName, user, password, host, tablePrefix string
	)
	dbType = configHelper.DBMysqlType
	dbName = configHelper.DBMysqlName
	user = configHelper.DBMysqlUser
	password = configHelper.DBMysqlPassword
	host = configHelper.DBMysqlHost
	tablePrefix = configHelper.DBMysqlTablePrefix

	connect := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err := gorm.Open(dbType, connect)
	if err != nil {
		log.Fatal(err)
	}
	// set Singular Table
	db.SingularTable(true)
	// 打印sql日志到控制台
	db.LogMode(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return fmt.Sprintf("%v%v", tablePrefix, defaultTableName)
	}
	SetDBConns(db)
	return db
}

// Set Max Open Conns
// Set MaxIdle Conns
func SetDBConns(db *gorm.DB) {
	MaxOpenConns, _ := strconv.Atoi(configHelper.MaxOpenConns)
	MaxIdleConns, _ := strconv.Atoi(configHelper.MaxIdleConns)
	db.DB().SetMaxOpenConns(MaxOpenConns)
	db.DB().SetMaxIdleConns(MaxIdleConns)
}

// 获取数据连接对象
func GetDB() *gorm.DB {
	return InitDB()
}

// 初始化mysql数据库并自动迁移模型
func InitModel() {
	DB.Init()
	db := GetDB()
	// 自动迁移模型
	db.AutoMigrate(&adminModel.Admin{})
}
