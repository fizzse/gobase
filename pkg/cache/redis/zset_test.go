package redis

import (
	"log"
	"testing"
)

func TestRedisClient_ZRangeWithScore(t *testing.T) {

}

func TestRedisClient_ZRangeByScoreWithScoreAndLimit(t *testing.T) {
	client := initTestClient()
	client.ZAdd("fizzse", 1, "c++")
	client.ZAdd("fizzse", 2, "go")
	client.ZAdd("fizzse", 3, "python")
	client.ZAdd("fizzse", 4, "java")
	client.ZAdd("fizzse", 4, "php")

	res, err := client.ZRangeByScoreWithScoreAndLimit("fizzse", 0, 4, 0, 2)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range res {
		log.Printf("key: %v,value: %v\n", k, v)
	}

	res, err = client.ZRangeByScoreWithScoreAndLimit("fizzse", 0, 4, 2, 3)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range res {
		log.Printf("key: %v,value: %v\n", k, v)
	}

	client.Del("fizzse")
}
