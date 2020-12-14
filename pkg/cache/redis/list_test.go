package redis

import (
	"log"
	"testing"
)

var (
	host     = "106.13.76.237"
	port     = "6379"
	password = "s"
)

func initTestClient() *RedisClient {
	testClient := NewClient(&RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
	})
	return testClient
}

func TestRedisClient_BRPop(t *testing.T) {
	client := initTestClient()
	//client.LPush("test", "c++")
	res, err := client.BLPop("test", 3)
	if err != nil {
		log.Println(res)
		log.Fatal(err)
	}

	log.Println(res)
}

func TestRedisClient_BLPop(t *testing.T) {
	client := initTestClient()
	client.LPush("test", "c++")
	client.LPush("test", "c#")
	res, err := client.BRPop("test", 3)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
