
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342672693268480>


//包优先级队列实现基于通道的优先级队列
//在任意类型上。它提供了一个
//一个自动操作循环，将一个函数应用于始终遵守的项
//他们的优先权。结构只是准一致的，即如果
//优先项是自动停止的，保证有一点
//当没有更高优先级的项目时，即不能保证
//有一点低优先级的项目存在
//但更高的不是

package priorityqueue

import (
	"context"
	"errors"
)

var (
	errContention  = errors.New("queue contention")
	errBadPriority = errors.New("bad priority")

	wakey = struct{}{}
)

//PriorityQueue是基本结构
type PriorityQueue struct {
	queues []chan interface{}
	wakeup chan struct{}
}

//New是PriorityQueue的构造函数
func New(n int, l int) *PriorityQueue {
	var queues = make([]chan interface{}, n)
	for i := range queues {
		queues[i] = make(chan interface{}, l)
	}
	return &PriorityQueue{
		queues: queues,
		wakeup: make(chan struct{}, 1),
	}
}

//运行是从队列中弹出项目的永久循环
func (pq *PriorityQueue) Run(ctx context.Context, f func(interface{})) {
	top := len(pq.queues) - 1
	p := top
READ:
	for {
		q := pq.queues[p]
		select {
		case <-ctx.Done():
			return
		case x := <-q:
			f(x)
			p = top
		default:
			if p > 0 {
				p--
				continue READ
			}
			p = top
			select {
			case <-ctx.Done():
				return
			case <-pq.wakeup:
			}
		}
	}
}

//push将项目推送到priority参数中指定的适当队列
//如果给定了上下文，它将一直等到推送该项或上下文中止为止。
//否则，如果队列已满，则返回errCompetition
func (pq *PriorityQueue) Push(ctx context.Context, x interface{}, p int) error {
	if p < 0 || p >= len(pq.queues) {
		return errBadPriority
	}
	if ctx == nil {
		select {
		case pq.queues[p] <- x:
		default:
			return errContention
		}
	} else {
		select {
		case pq.queues[p] <- x:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	select {
	case pq.wakeup <- wakey:
	default:
	}
	return nil
}

