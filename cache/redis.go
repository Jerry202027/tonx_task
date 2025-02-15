package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	// global redis client
	Rdb *redis.Client
	// for redis
	Ctx = context.Background()
)

func InitRedis(addr, password string, db int) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func DecrementRemainingSeats(flightID uint, defaultValue int) (int, error) {
	key := fmt.Sprintf("flight:%d:remaining", flightID)
	// use Lua script to manipulate RemainingSeats
	script := redis.NewScript(`
	local seats = tonumber(redis.call("GET", KEYS[1]))
	if not seats then
		seats = tonumber(ARGV[1])
		redis.call("SET", KEYS[1], seats)
	end
	if seats <= 0 then
		return -1
	end
	seats = seats - 1
	redis.call("SET", KEYS[1], seats)
	return seats
`)
	// defaultStr := fmt.Sprintf("%d", defaultValue)
	result, err := script.Run(Ctx, Rdb, []string{key}, defaultValue).Result()
	if err != nil {
		return 0, err
	}

	if seats, ok := result.(int64); ok {
		return int(seats), nil
	}

	return 0, fmt.Errorf("unexpected result type")
}
