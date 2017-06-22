package main

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

var (
	// 定义常量
	RedisClient     *redis.Pool
	REDIS_HOST string
	REDIS_DB   int
)

func init() {
	// 从配置文件获取redis的ip以及db
	REDIS_HOST = "10.0.30.120:6379"
	REDIS_DB = 1
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     1,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}
func main(){
	// 从池里获取连接
	c := RedisClient.Get()

	r,err := c.Do("PING")
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println(r)
	}

	// 用完后将连接放回连接池
	defer c.Close()
}
