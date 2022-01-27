package kafka

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestPublisher(t *testing.T) {
	var (
		testTopic   = "sample-topic"
		testBrokers = []string{"127.0.0.1:9092"}
	)

	pub := NewPublisher(testBrokers)
	defer pub.Close()

	log.Println(time.Now())

	data := `{
	"task_id": "97f3a21a-92cd-11eb-b89c-02423cbe895b",
	"product": 3,
	"mac": "582D34006822",
	"version": "1.0.3_0006",
	"timestamp": 1617698123
}`

	for i := 0; i < 1; i++ {
		err := pub.Publish(context.Background(), Event{
			Topic:      testTopic,
			Properties: map[string]string{"traceId": "11"},
			Payload:    []byte(data),
		})
		if err != nil {
			log.Fatal("send message failed", err)
		}
	}

	log.Println(time.Now())
}
