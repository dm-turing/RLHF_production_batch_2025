package lib

import (
	"context"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func InitRedis(addr string) redis.ConnectionPool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}
		c.SetDeadline(time.Now().Add(time.Second * 3))
		return c, nil
	}, 10) // 10 is the maximum number of active connections to Redis
}

func Get(ctx context.Context, key string) ([]byte, error) {
	c := pool.Get()
	defer c.Close()
	c.SetContext(ctx)

	val, err := c.Do("GET", key)
	if err != nil {
		return nil, err
	}
	b, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("could not convert value to []byte")
	}
	return b, nil
}

func Set(ctx context.Context, key string, val interface{}) error {
	c := pool.Get()
	defer c.Close()
	c.SetContext(ctx)

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	_, err = c.Do("SET", key, b)
	return err
}

func Del(ctx context.Context, key string) (int64, error) {
	c := pool.Get()
	defer c.Close()
	c.SetContext(ctx)

	return c.Do("DEL", key)
}
