package biz

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/fizzse/gobase/pkg/mq/kafka"
)

const (
	namespace = "gobase" // serverName
)

var (
	msgCounter     *prometheus.CounterVec
	msgCounterOnce sync.Once
)

const (
	statusOk     = "ok"
	statusFailed = "failed"
)

func NewMetricCounter(name, help string, labels ...string) *metricCounter {
	msgCounterOnce.Do(func() { // init in first call
		msgCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      name, //   "process_msg_count",
				Help:      help, // "Total number of message of iot device",
			}, labels,
		)
	})

	return &metricCounter{labels: labels, status: statusOk}
}

type metricCounter struct {
	labels []string
	status string
}

func (c *metricCounter) metric() {
	lvs := c.labels
	msgCounter.WithLabelValues(lvs...).Inc()
}

// DealMsg kafka

type ExampleMsg struct {
	Key string `json:"key"`
}

func (b *SampleBiz) DealMsg(ctx context.Context, msg kafka.Event) (err error) { // 回调函数
	parsedMsg := ExampleMsg{}
	err = json.Unmarshal(msg.Payload, &parsedMsg)
	if err != nil {
		return err
	}

	m := NewMetricCounter("process_msg_count", "Total number of message of iot device", parsedMsg.Key)

	defer func() {
		if err != nil {
			m.status = statusOk
		}

		m.metric()
	}()
	return
}
