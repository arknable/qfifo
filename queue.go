package qfifo

import "sync"

// Queue is a queue with FIFO behavior.
type Queue struct {
	list []interface{}
	lock sync.Mutex
}

// QueueOptions is optional settings to be used when creating a queue.
type QueueOptions struct {
	InitialSize int
}

// Push adds new element into the end of queue.
func (q *Queue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.list = append(q.list, v)
}

// Pop removes element from the end of queue.
func (q *Queue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	length := len(q.list)
	if length == 0 {
		return nil
	}

	lastIdx := length - 1
	v := q.list[lastIdx]
	q.list[lastIdx] = nil
	q.list = q.list[:lastIdx]
	return v
}

// Clear removes stored elements
func (q *Queue) Clear() {
	q.lock.Lock()
	defer q.lock.Unlock()

	for i := range q.list {
		q.list[i] = nil
	}
	q.list = q.list[:0]
}

// IsEmpty returns true if queue has no element stored, otherwise false.
func (q *Queue) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.list) == 0
}

// Default slice size
const defaultSize = 10

func New(opts *QueueOptions) *Queue {
	if opts == nil {
		opts = &QueueOptions{
			InitialSize: defaultSize,
		}
	}
	return &Queue{
		list: make([]interface{}, 0, opts.InitialSize),
	}
}
