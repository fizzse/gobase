package redis

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var (
	target    = "15156133965"
	category  = "sms"
	keyFormat = "qingPush_%s_%s"
)

const count = 10

func TestCounter_ZINCR(t *testing.T) {
	client := initTestClient()
	counter := NewCounter(60, client)

	day := time.Now().Format("2006-01-02")
	key := fmt.Sprintf(keyFormat, category, day)

	for i := 1; i <= count; i++ {
		num, err := counter.ZINCR(key, target)
		if err != nil {
			log.Fatal(err)
		}

		if num != i {
			log.Fatalf("num: %d must be %d", num, i)
		}
	}
	// client.Del(key)
	log.Printf("test ok")
}

func TestCounter_GetZScore(t *testing.T) {
	client := initTestClient()
	counter := NewCounter(60, client)
	day := time.Now().Format("2006-01-02")
	key := fmt.Sprintf(keyFormat, category, day)

	num, err := counter.GetZScore(key, target)
	if err != nil {
		log.Fatal(err)
	}
	client.Del(key)
	if num != count {
		log.Fatalf("num: %d must be %d", num, count)
	}
	log.Printf("num: %d, count: %d", num, count)
	log.Printf("test ok")
}
