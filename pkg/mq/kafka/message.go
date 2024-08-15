package kafka

import "context"

type Config struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupId string   `yaml:"groupId"`
}

type Event struct {
	// Key sets the key of the message for routing policy
	Key string
	// Payload for the message
	Payload []byte
	// Properties attach application defined properties on the message
	Properties map[string]string
	// topic
	Topic string
	// offset
	Offset int64
	// 分区
	Partition int
}

type Handler func(context.Context, Event) error
