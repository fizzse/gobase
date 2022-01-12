package biz

import (
	"context"
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/fizzse/gobase/pkg/mq/kafka"
)

const (
	namespace = "gobase" // serverName
)

var (
	labels      = []string{"drive", "productKey", "status"}
	MegCountReq = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "process_msg_count",
			Help:      "Total number of message of iot device",
		}, labels,
	)
)

const (
	statusOk     = "ok"
	statusFailed = "failed"
)

type metricCounter struct {
	labels []string
	status string
}

func (c *metricCounter) metric() {
	lvs := c.labels
	MegCountReq.WithLabelValues(lvs...).Inc()
}

func NewMetricCounter(labels ...string) *metricCounter {
	return &metricCounter{labels: labels, status: statusOk}
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

	m := NewMetricCounter(parsedMsg.Key)

	defer func() {
		if err != nil {
			m.status = statusOk
		}

		m.metric()
	}()
	return
}
