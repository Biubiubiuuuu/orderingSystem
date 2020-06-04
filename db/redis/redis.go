package redis

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/gomodule/redigo/redis"
)

type RedisDataBase struct {
	Redis *redis.Pool
}

var DB *RedisDataBase
var once sync.Once

// 初始化redis连接
func (db *RedisDataBase) InitConn() {
	once.Do(func() {
		DB = &RedisDataBase{
			Redis: InitRedisConePool(),
		}
	})
}

// 初始化redis连接池
func InitRedisConePool() *redis.Pool {
	host := configHelper.DBRedisHost
	db, _ := strconv.Atoi(configHelper.DBRedisDb)
	pass := configHelper.DBRedisPassword
	maxActive, _ := strconv.Atoi(configHelper.DBRedisMaxActive)
	maxIdle, _ := strconv.Atoi(configHelper.DBRedisMaxIdle)
	idleTimeout, _ := strconv.ParseInt(configHelper.DBRedisIdleTimeout, 10, 64)
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host, redis.DialDatabase(db), redis.DialConnectTimeout(time.Duration(idleTimeout)*time.Second))
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			if pass != "" {
				if _, err := c.Do("AUTH", pass); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// 获取redis连接池
func GetRedisDB() redis.Conn {
	return InitRedisConePool().Get()
}

// 关闭redis
func CloseRedis() {
	DB.Redis.Close()
}
