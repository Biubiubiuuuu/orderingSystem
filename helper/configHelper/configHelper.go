package configHelper

import (
	"log"

	"github.com/go-ini/ini"
)

var (
	Cfg                *ini.File
	RunMode            string
	HTTPPort           string
	JwtSecret          string
	JwtName            string
	DBMysqlType        string
	DBMysqlName        string
	DBMysqlUser        string
	DBMysqlPassword    string
	DBMysqlHost        string
	DBMysqlTablePrefix string
	Version            string
	Static             string
	LogDir             string
	ImageDir           string
	MaxIdleConns       string
	MaxOpenConns       string
	DBRedisHost        string
	DBRedisPassword    string
	DBRedisDb          string
	DBRedisMaxActive   string
	DBRedisMaxIdle     string
	DBRedisIdleTimeout string
)

// 初始化配置信息
func init() {
	var err error
	Cfg, err = ini.Load("./config/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'config/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadMysql()
	LoadRedis()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	HTTPPort = sec.Key("HTTP_PORT").MustString("8060")
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	JwtName = sec.Key("JWT_NAME").MustString("token")
	Version = sec.Key("VERSION").MustString("1.0")
	Static = sec.Key("STATIC").MustString("static/")
	LogDir = sec.Key("LOGDIR").MustString("static/log/")
	ImageDir = sec.Key("IMAGEDIR").MustString("static/image/")
}

func LoadMysql() {
	sec, err := Cfg.GetSection("mysql")
	if err != nil {
		log.Fatalf("Fail to get section 'mysql': %v", err)
	}
	DBMysqlType = sec.Key("TYPE").MustString("mysql")
	DBMysqlName = sec.Key("NAME").MustString("test")
	DBMysqlUser = sec.Key("USER").MustString("root")
	DBMysqlPassword = sec.Key("PASSWORD").MustString("")
	DBMysqlHost = sec.Key("HOST").MustString("127.0.0.1:3306")
	DBMysqlTablePrefix = sec.Key("TABLE_PREFIX").MustString("")
	MaxIdleConns = sec.Key("MAXIDLECONNS").MustString("20")
	MaxOpenConns = sec.Key("MAXOPENCONNS").MustString("20")
}

func LoadRedis() {
	sec, err := Cfg.GetSection("redis")
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}
	DBRedisHost = sec.Key("HOST").MustString("127.0.0.1:3306")
	DBRedisPassword = sec.Key("PASSWORD").MustString("")
	DBRedisDb = sec.Key("DB").MustString("0")
	DBRedisMaxActive = sec.Key("MAXACTIVE").MustString("20")
	DBRedisMaxIdle = sec.Key("MAXIDLE").MustString("")
	DBRedisIdleTimeout = sec.Key("IDLETIMEOUT").MustString("0")
}
