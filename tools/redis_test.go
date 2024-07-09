package tools

import (
	"context"
	"testing"
	"time"
	"tinyIM/config"
)

func TestRedisConnect(t *testing.T) {
	redisOpt := RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}

	client := GetRedisInstance(redisOpt)

	// 尝试 ping Redis 服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %s", err.Error())
	}

	if pong != "PONG" {
		t.Fatalf("Unexpected response from Redis: %s", pong)
	}

	t.Log("Redis connection successful")
}
