package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	pool *redis.Pool
	redisHost = "127.0.0.1:6379"
)

//newRedisPool  创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:50,
		MaxActive:30,
		IdleTimeout:300*time.Second,
		Dial: func() (redis.Conn, error) {
			//打开链接
			conn, err := redis.Dial("tcp", redisHost)
			if err !=nil {
				return nil, err
			}
			return conn,nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}
