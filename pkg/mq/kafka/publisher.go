package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

// ErrEventFull is a message event chan full.
var ErrEventFull = errors.New("message event chan full")

const (
	// RequireNone the producer won’t even wait for a response from the broker.
	RequireNone kafka.RequiredAcks = kafka.RequireNone
	// RequireOne the producer will consider the write successful when the leader receives the record.
	RequireOne kafka.RequiredAcks = kafka.RequireOne
	// RequireAll the producer will consider the write successful when all of the in-sync replicas receive the record.
	RequireAll kafka.RequiredAcks = kafka.RequireAll
)

// PublisherOption is a publisher options.
type PublisherOption func(*Publisher)

// ReadTimeout with read timeout option.
func ReadTimeout(d time.Duration) PublisherOption {
	return func(o *Publisher) {
		o.readTimeout = d
	}
}

// WriteTimeout with write timeout option.
func WriteTimeout(d time.Duration) PublisherOption {
	return func(o *Publisher) {
		o.writeTimeout = d
	}
}

// EventBuffer with event buffer option.
func EventBuffer(n int) PublisherOption {
	return func(o *Publisher) {
		o.eventBuffer = n
	}
}

// RequiredAcks with required acks option.
func RequiredAcks(acks kafka.RequiredAcks) PublisherOption {
	return func(o *Publisher) {
		o.requiredAcks = acks
	}
}

func BatchSize(size int) PublisherOption {
	return func(o *Publisher) {
		o.batchSize = size
	}
}

type pubEvent struct {
	ctx      context.Context
	event    Event
	callback func(err error)
}

type Publisher struct {
	brokers      []string
	readTimeout  time.Duration
	writeTimeout time.Duration
	eventBuffer  int
	requiredAcks kafka.RequiredAcks
	writer       *kafka.Writer
	eventChan    chan pubEvent
	batchSize    int           // 批量处理
	batchTimeOut time.Duration // 批量处理
}

// NewPublisher new a kafka publisher.
func NewPublisher(brokers []string, opts ...PublisherOption) *Publisher {
	pub := &Publisher{
		brokers:      brokers,
		readTimeout:  500 * time.Millisecond,
		writeTimeout: 500 * time.Millisecond,
		batchSize:    1,
		batchTimeOut: 500 * time.Millisecond,
		eventBuffer:  1000,
		requiredAcks: RequireOne,
	}
	for _, o := range opts {
		o(pub)
	}
	pub.writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		RequiredAcks: pub.requiredAcks,
		ReadTimeout:  pub.readTimeout,
		WriteTimeout: pub.writeTimeout,
		BatchSize:    pub.batchSize,
		BatchTimeout: pub.batchTimeOut,
	}

	return pub
}

func (p *Publisher) Publish(ctx context.Context, event Event) error {
	headers := make([]kafka.Header, 0, len(event.Properties))
	for k, v := range event.Properties {
		headers = append(headers, kafka.Header{Key: k, Value: []byte(v)})
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic:   event.Topic,
		Key:     []byte(event.Key),
		Value:   event.Payload,
		Headers: headers,
	})
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
