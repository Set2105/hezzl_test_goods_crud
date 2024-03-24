package redis

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type RedisSettings struct {
	Addr     string `xml:"redisAddr"`
	Port     string `xml:"redisPort"`
	User     string `xml:"redisUser"`
	Password string `xml:"redisPassword"`
	Db       int    `xml:"redisDb"`
}

func (settings *RedisSettings) Valid() error {
	if settings.Addr == "" {
		settings.Addr = "redis"
	}
	if settings.Port == "" {
		settings.Port = "6379"
	}
	return nil
}

type Redis struct {
	C *redis.Client
}

func (r *Redis) GetByte(key string) ([]byte, error) {
	b, err := r.C.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("Redis.GetByte: %s", err)
	}
	return b, nil
}

func (r *Redis) ChacheByte(key string, data []byte, cTime time.Duration) error {
	if _, err := r.C.Set(context.Background(), key, data, cTime).Result(); err != nil {
		return fmt.Errorf("Redis.ChacheByte.%s", err)
	}
	return nil
}

func InitRedis(s *RedisSettings) (c *Redis, err error) {
	err = s.Valid()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     s.Addr + ":" + s.Port,
		Password: s.Password,
		DB:       s.Db,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{C: rdb}, nil
}
