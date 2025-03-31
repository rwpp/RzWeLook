package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/rwpp/RzWeLook/config"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
