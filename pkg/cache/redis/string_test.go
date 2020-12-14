package redis

import (
	"log"
	"testing"
)

func TestRedisClient_INCR(t *testing.T) {
	client := initTestClient()
	key := "1"

	for i := 1; i < 5; i++ {
		counter, err := client.INCR(key)
		if err != nil {
			log.Fatal(err)
		}

		if counter != i {
			log.Fatalf("counter: %d must be %d", counter, i)
		}
	}

	client.Del(key)
	log.Printf("test ok")
}

func TestRedisClient_Keys(t *testing.T) {
	client := initTestClient()

	res, err := client.Keys("*")
	if err != nil {
		log.Fatal(err)
	}

	for _, key := range res {
		log.Println(key)
	}
}
