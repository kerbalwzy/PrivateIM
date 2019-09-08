package utils

import (
	"github.com/go-redis/redis"
	"log"
)

func GetRedisClient(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // default no password set
		DB:       db,       // use default DB
	})

	_, err := client.Ping().Result()
	if nil != err {
		log.Fatal(err)
	}
	return client
}
