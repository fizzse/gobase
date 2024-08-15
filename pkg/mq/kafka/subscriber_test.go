package kafka

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestSubscriber(t *testing.T) {
	topic := "hello"
	groupId := "world"
	broker := []string{"127.0.0.1:9200"}

	r := NewSubscriber(&Config{Brokers: broker, Topic: topic, GroupId: groupId})

	go func() {
		err := r.Subscribe(context.Background(), func(ctx context.Context, event Event) error {
			log.Printf("sub: key=%s value=%s header=%v", event.Key, event.Payload, event.Properties)
			return nil
		})

		if err != nil {
			log.Fatal("subscribe failed ", err)
		}
	}()

	time.Sleep(5 * time.Minute)
	r.Close()
	time.Sleep(time.Second)
}

func TestSubscriber_ReadMessage(t *testing.T) {
	topic := "hello"
	groupId := "world"
	broker := []string{"127.0.0.1:9200"}

	r := NewSubscriber(&Config{Brokers: broker, Topic: topic, GroupId: groupId})

	go func() {
		for {
			message, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Fatal("subscribe failed ", err)
			}

			log.Printf("sub: key=%s value=%s header=%v", message.Key, message.Value, message.Headers)
		}

	}()

	time.Sleep(10 * time.Second)
	r.Close()
	time.Sleep(time.Second)
}
