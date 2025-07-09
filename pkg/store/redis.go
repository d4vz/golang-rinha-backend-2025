package store

import (
	"context"
	"fmt"

	"github.com/d4vz/rinha-de-backend-2025/config"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func ConnectRedis() error {
	redisHost := config.GetEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := config.GetEnvOrDefaultInt("REDIS_PORT", 6379)
	connString := fmt.Sprintf("%s:%d", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr: connString,
	})

	_, err := RedisClient.Ping(ctx).Result()

	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w connection string: %s", err, connString)
	}

	return nil
}

func GetRedis() *redis.Client {
	return RedisClient
}
