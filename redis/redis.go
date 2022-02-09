package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

// Start 启动redis链接
func Start(url string, pwd string) {
	pool = newPool(url, pwd)

	// ping
	for {
		conn := pool.Get()
		if conn.Err() == nil && conn != nil {
			log.Printf("redis =>[%s] established.", url)
			conn.Close()
			break
		} else {
			log.Println("retry 2 seconds later ...")
			time.Sleep(time.Duration(2) * time.Second)
		}
	}
}

// Get 获取redis链接
func Get() redis.Conn {
	return pool.Get()
}

func newPool(addr string, pwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		MaxActive:   100,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			log.Println("open redis...")
			con, err := redis.Dial("tcp", addr)
			if err != nil {
				fmt.Println(err)
			} else if len(pwd) != 0 {
				_, err = con.Do("AUTH", pwd)
				if err != nil {
					fmt.Println(err)
				}
			}

			return con, err
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
