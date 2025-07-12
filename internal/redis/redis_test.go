package Redis

import (
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestInitRedis(T *testing.T) {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		T.Errorf("Redis is not connected")
	}
}
