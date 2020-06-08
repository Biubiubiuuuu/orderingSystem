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
	Redis redis.Conn
}

var DB *RedisDataBase
var once sync.Once

// 初始化redis连接
func (db *RedisDataBase) InitConn() {
	once.Do(func() {
		DB = &RedisDataBase{
			Redis: InitRedisConnPool(),
		}
	})
}

// 初始化redis连接池
func InitRedisConnPool() redis.Conn {
	host := configHelper.DBRedisHost
	db, _ := strconv.Atoi(configHelper.DBRedisDb)
	pass := configHelper.DBRedisPassword
	maxActive, _ := strconv.Atoi(configHelper.DBRedisMaxActive)
	maxIdle, _ := strconv.Atoi(configHelper.DBRedisMaxIdle)
	idleTimeout, _ := strconv.ParseInt(configHelper.DBRedisIdleTimeout, 10, 64)
	pool := redis.Pool{
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
					log.Fatal(err)
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
			log.Fatal(err)
			return err
		},
	}
	return pool.Get()
}

// 获取redis连接池
func GetRedisDB() redis.Conn {
	return InitRedisConnPool()
}
