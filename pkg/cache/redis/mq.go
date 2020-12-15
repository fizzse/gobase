package redis

/*
 * redis消息队列
 */

import (
	"context"
	"log"
)

const (
	statusAlive = 1 // 工作状态
	statusStop  = 2 // 退出状态

	DefaultTimeout   = 1            // 获取不到数据休眠时间
	DefaultStoreTime = 24 * 60 * 60 // 数据持久化时间
)

// 消息处理函数类型
type HandelFunc func(context.Context, string) error

type MessageQueue struct {
	status       int                // 运行状态
	workerCount  int                // 协程数
	recvTimeout  int                // 获取数据超时时间
	msgStoreTime int                // 消息持久化时间
	queueName    string             // 名称 唯一标识
	entity       *Client            // redis 实例
	handelFunc   HandelFunc         // 处理消息函数
	cancelCtx    context.Context    // 退出控制
	cancelFunc   context.CancelFunc // 退出控制
	debugModel   bool               // 调试模式 log打印日志
}

func NewQueue(name string, workerCount int, handelFunc HandelFunc, entity *Client) *MessageQueue {
	q := &MessageQueue{queueName: name, workerCount: workerCount, handelFunc: handelFunc, entity: entity}
	if q.recvTimeout == 0 {
		q.recvTimeout = DefaultTimeout
	}

	if q.msgStoreTime == 0 {
		q.msgStoreTime = DefaultStoreTime
	}

	q.cancelCtx, q.cancelFunc = context.WithCancel(context.Background())
	return q
}

func (q *MessageQueue) Debug(ok bool) {
	q.debugModel = ok
}

func (q *MessageQueue) Push(msg string) error {
	return q.entity.LPush(q.queueName, msg)
}

func (q *MessageQueue) Pop() string {
	msg, err := q.entity.BRPop(q.queueName, q.recvTimeout)
	if err != nil {
		return ""
	}

	return msg
}

func (q *MessageQueue) Run() (func(), error) {
	q.status = statusAlive
	go func() {
		for i := 0; i < q.workerCount; i++ {
			go func(id int, exitChan <-chan struct{}) {
				if q.debugModel {
					log.Printf("%s queue worker: %d ready to work\n", q.queueName, id)
				}

				for {
					select {
					case <-exitChan:
						if q.debugModel {
							log.Printf("%s queue worker: %d exit\n", q.queueName, id)
						}
						return

					default:
						pushMsg := q.Pop()
						if pushMsg == "" {
							if q.debugModel {
								log.Printf("%s queue worker: %d no message!\n", q.queueName, id)
							}

							continue
						}

						if q.debugModel {
							log.Printf("%s queue worker: %d deal msg: %s\n", q.queueName, id, pushMsg)
						}

						// 处理数据 FIXME ctx
						q.handelFunc(context.Background(), pushMsg)
					}
				}
			}(i, q.cancelCtx.Done())
		}
	}()

	return q.Stop, nil
}

// 退出的时候 给key加过期时间
func (q *MessageQueue) Stop() {
	q.status = statusStop
	//q.entity.Expire(q.queueName, q.msgStoreTime)
	q.cancelFunc()
}

func (q *MessageQueue) Running() bool {
	return q.status == statusAlive
}
