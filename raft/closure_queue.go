package raft

import (
	"errors"
	"github.com/AllenShaw19/raft/log"
	"sync"
)

type ClosureQueue struct {
	mutex      sync.Mutex
	firstIndex int64
	queue      []Closure // TODO: 可以换成一个无锁队列
}

func (q *ClosureQueue) Clear() {
	q.mutex.Lock()
	savedQueue := q.queue
	q.queue = make([]Closure, 0)
	q.firstIndex = 0
	q.mutex.Unlock()

	for _, c := range savedQueue {
		c.GetStatus().SetError(int32(RaftError_EPERM), "leader stepped down")
		go c.Run()
	}
}

func (q *ClosureQueue) ResetFirstIndex(firstIndex int64) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.queue) != 0 {
		log.Fatal("queue is nil")
	}
	q.firstIndex = firstIndex
}

func (q *ClosureQueue) AppendPendingClosure(c Closure) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.queue = append(q.queue, c)
}

func (q *ClosureQueue) PopClosureUntil(index int64) (out []Closure, outFirstIndex int64, err error) {
	out = make([]Closure, 0)
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.queue) == 0 || index < q.firstIndex {
		return nil, index + 1, nil
	}

	if index > q.firstIndex+int64(len(q.queue))-1 {
		log.Error("invalid index=%d, first_index=%d, closure_queue_size=%d", index, q.firstIndex, len(q.queue))
		return nil, 0, errors.New("invalid index")
	}
	outFirstIndex = q.firstIndex
	for i := q.firstIndex; i <= index; i++ {
		out = append(out, q.queue[i])
	}
	q.queue = q.queue[index-q.firstIndex:]
	q.firstIndex = index + 1
	return out, outFirstIndex, nil
}
