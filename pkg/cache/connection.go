package cache

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func GetRedisConnection() (*redis.Client, error) {
	address := fmt.Sprintf("%v:%v", viper.GetString("redis.host"), viper.GetString("redis.port"))

	redisClient := redis.NewClient(&redis.Options{
		Addr: address,
	})
	
	if redisClient == nil {
		return nil, errors.New("failed connect to redis")
	}

	return redisClient, nil
}