package cache

/*
this file will manage the caching
using redis

imports required:
redis
the utils package

functions:
func initCache() *redis.client
func get(key string) (interface{}, bool) get from cache
func set(key string, value interface{}) set in cache

here interface{} means that it can hold any type of data
*/
import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrCacheMiss = errors.New("cache: key not found")

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
	// the password is redis level auth, this is mandatory for communication
	// with redis

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return rdb
}

// shld run later on a gin handlerfunc to add to list of calls
func SetToRedis(rdb *redis.Client, ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return rdb.Set(ctx, key, value, ttl).Err()
}

// shld run later on a gin handlerfunc to add to list of calls
func GetFromRedis(rdb *redis.Client, ctx context.Context, key string) ([]byte, bool, error) {
	rep, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key doesnt exist")
		return nil, false, ErrCacheMiss
	} else if err != nil {
		return nil, false, err
	}

	return []byte(rep), true, nil
}