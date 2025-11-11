package cache

/*
this file will manage the caching
using ristretto cacher

imports required:
dgraph-io/ristretto
the utils package

functions:
func initCache() *ristretto.Cache
func get(key string) (interface{}, bool) get from cache
func set(key string, value interface{}) set in cache

here interface{} means that it can hold any type of data
*/
import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func ConnectToRedis() *redis.Client {
	// NewClient returns a client to the 
	// Redis Server specified by Options.
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0, // this is an internal redis DB no
		Protocol: 2,
	})

	// redis has multiple logical DBs, from 0 to 15
	// we will here use the standard, ie 0
	// Protocol: 2 -> specifies the Redis Serialization Protocol (RESP)
	// 2 is the older version (stable), 3 is new, intro: redis 6
	// the password is redis leve auth, this is mandatory for communication
	// with redis

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return rdb
}

// shld run later on a gin handlerfunc to add to list of calls
func SetToRedis(rdb *redis.Client, ctx context.Context, key string, value interface{}) {
	if err := rdb.Set(ctx, key, value, 0).Err(); err != nil {
		panic(err)
	}
}

// shld run later on a gin handlerfunc to add to list of calls
func GetFromRedis(rdb *redis.Client, ctx context.Context, key string) string {
	rep, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key doesnt exist")
		return ""
	} else if err != nil {
		panic(err)
	}

	return rep
}