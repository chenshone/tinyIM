package tools

import (
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var RedisClientMap = map[string]*redis.Client{}
var syncLock sync.Mutex

type RedisOption struct {
	Address  string
	Password string
	Db       int
}

func GetRedisInstance(redisOpt RedisOption) *redis.Client {
	addr := redisOpt.Address
	syncLock.Lock()
	defer syncLock.Unlock()

	// 如果客户端实例已存在，直接返回
	if client, ok := RedisClientMap[addr]; ok {
		return client
	}

	// 初始化 Redis 客户端实例
	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        redisOpt.Password,
		DB:              redisOpt.Db,
		ConnMaxLifetime: 20 * time.Second,
	})

	// 将新的客户端实例存储在映射中
	RedisClientMap[addr] = client
	return client
}
