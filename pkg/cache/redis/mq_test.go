package redis

import (
	"context"
	"log"
	"testing"
	"time"
)

func testHandelFunc(ctx context.Context, msg string) error {
	log.Println(msg)
	return nil
}

func TestMessageQueue(t *testing.T) {
	entity := initTestClient()
	mq := NewQueue("test-queue", 2, testHandelFunc, entity)
	mq.Debug(true)
	//mq.debugModel = true

	close, err := mq.Run()
	if err != nil {
		log.Fatal(err)
	}

	//mq.Push("C++")
	//mq.Push("Go")

	time.Sleep(4 * time.Second)
	//mq.Push("JAVA")
	//mq.Push("PYTHON")

	close()

	time.Sleep(5 * time.Second)
}
