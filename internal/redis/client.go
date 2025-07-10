package Redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

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

	expireTime := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	ttl := time.Until(expireTime)

	err := Client.HSet(ctx, "user:"+username,
		"calories", limits.Calories,
		"proteins", limits.Proteins,
		"fats", limits.Fats,
		"carbs", limits.Carbs,
	).Err()

	if err != nil {
		log.Printf("error set limits %v", err)
		return err
	}

	return Client.Expire(ctx, "user:"+username, ttl).Err()
}

func GetUserLimits(username string) (limits.DailyLimits, error) {
	result, err := Client.HGetAll(ctx, "user:"+username).Result()
	if err != nil {
		return limits.DailyLimits{}, err
	}

	var (
		calories, proteins, fats, carbs int
	)

	if val, err := strconv.Atoi(result["calories"]); err == nil {
		calories = val
	}
	if val, err := strconv.Atoi(result["proteins"]); err == nil {
		proteins = val
	}
	if val, err := strconv.Atoi(result["fats"]); err == nil {
		fats = val
	}
	if val, err := strconv.Atoi(result["carbs"]); err == nil {
		carbs = val
	}

	return limits.DailyLimits{
		Calories: calories,
		Proteins: proteins,
		Fats:     fats,
		Carbs:    carbs,
	}, nil
}

func SubtractMeal(username string, calories, proteins, fats, carbs int) (limits.DailyLimits, error) {

	_, err := Client.HIncrBy(ctx, "user:"+username, "calories", -int64(calories)).Result()
	if err != nil {
		return limits.DailyLimits{}, err
	}
	_, err = Client.HIncrBy(ctx, "user:"+username, "proteins", -int64(proteins)).Result()
	if err != nil {
		return limits.DailyLimits{}, err
	}
	_, err = Client.HIncrBy(ctx, "user:"+username, "fats", -int64(fats)).Result()
	if err != nil {
		return limits.DailyLimits{}, err
	}
	_, err = Client.HIncrBy(ctx, "user:"+username, "carbs", -int64(carbs)).Result()
	if err != nil {
		return limits.DailyLimits{}, err
	}

	return GetUserLimits(username)
}
