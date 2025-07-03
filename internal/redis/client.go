package Redis

import (
	"context"
	"os"
	"strconv"

	limits "telegrambot/internal/limits"

	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

func InitRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := Client.Ping(ctx).Result()
	return err
}

func SetUserLimits(username string, limits limits.DailyLimits) error {

	data := map[string]int{
		"calories": limits.Calories,
		"proteins": limits.Proteins,
		"fats":     limits.Fats,
		"carbs":    limits.Carbs,
	}

	err := Client.HSet(ctx, "user:"+username, data).Err()
	if err != nil {
		return err
	}

	return err
}

func GetUserLimits(username string) (limits.DailyLimits, error) {
	result, err := Client.HGetAll(ctx, "user:"+username).Result()
	if err != nil {
		return limits.DailyLimits{}, err
	}

	calories, _ := strconv.Atoi(result["calories"])
	proteins, _ := strconv.Atoi(result["proteins"])
	fats, _ := strconv.Atoi(result["fats"])
	carbs, _ := strconv.Atoi(result["carbs"])

	return limits.DailyLimits{
		Calories: calories,
		Proteins: proteins,
		Fats:     fats,
		Carbs:    carbs,
	}, nil
}
