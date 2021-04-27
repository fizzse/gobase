package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Subscriber struct {
	reader  *kafka.Reader
	groupId string
}

// SubscriberOption is a subscriber option.
type SubscriberOption func(*Subscriber)

func ConsumerGroup(id string) SubscriberOption {
	return func(o *Subscriber) {
		o.groupId = id
	}
}

// NewSubscriber new a kafka subscriber.
func NewSubscriber(topic string, brokers []string, opts ...SubscriberOption) *Subscriber {
	sub := &Subscriber{}
	for _, o := range opts {
		o(sub)
	}

	dialer := kafka.DefaultDialer
	dialer.Timeout = 2 * time.Second

	sub.reader = kafka.NewReader(kafka.ReaderConfig{
		Topic:   topic,
		GroupID: sub.groupId,
		Brokers: brokers,
		MaxWait: 2 * time.Second,
		Dialer:  dialer,
	})
	return sub
}

func (s *Subscriber) Subscribe(ctx context.Context, h Handler) error {
	for {
		msg, err := s.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		header := make(map[string]string, len(msg.Headers))
		for _, h := range msg.Headers {
			header[h.Key] = string(h.Value)
		}

		_ = h(context.Background(), Event{
			Key:        string(msg.Key),
			Payload:    msg.Value,
			Properties: header,
			Offset:     msg.Offset,
			Partition:  msg.Partition,
		})
		_ = s.reader.CommitMessages(ctx, msg)
	}
}

func (s *Subscriber) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return s.reader.ReadMessage(ctx)
}

func (s *Subscriber) GetEntity() *kafka.Reader {
	return s.reader
}

func (s *Subscriber) Close() error {
	return s.reader.Close()
}
