package redis_example

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func Connect() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.9.227:6379",
		Password: "123456",
		DB:       0,
	})
	err := rdb.Set(ctx, "name", "caiyanan", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	rdb.Close()
}
