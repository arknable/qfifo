package qfifo

import (
	"errors"
	"sync"
	"time"
)

// ErrPublishFunctionUnset occurs if NewPublisher() called with PublishFunc unset.
var ErrPublishFunctionUnset = errors.New("publish function is mandatory")

// PublisherArgs is collection of settings to create a publisher.
type PublisherArgs struct {
	// Queue is the queue to be used. If not set then a queue with default settings
	// will be created.
	Queue *Queue

	// SleepInterval is the sleep duration of internal go routine
	// that publish queued elements, default value is 100 ms.
	SleepInterval time.Duration

	// PublishFunc is the function to receive published element.
	// If unset, NewPublisher() returns ErrPublishFunctionUnset.
	PublishFunc func(*Publisher, interface{})
}

// Publisher extends Queue with self-publishing feature.
// Self-publishing means it automatically pops element one at a time
// and sends that popped element to a publish function, i.e. PublishFunc.
type Publisher struct {
	queue         *Queue
	sleepInterval time.Duration
	publishFunc   func(*Publisher, interface{})

	keepPublishing    bool
	publishingStarted bool
	wait              sync.WaitGroup
	waitStart         sync.WaitGroup
}

// Push adds new element into the end of queue.
func (q *Publisher) Push(v interface{}) {
	q.queue.Push(v)
}

// Pop removes element from the end of queue.
func (q *Publisher) Pop() interface{} {
	return q.queue.Pop()
}

// Clear removes stored elements
func (q *Publisher) Clear() {
	q.queue.Clear()
}

// IsEmpty returns true if queue has no element stored, otherwise false.
func (q *Publisher) IsEmpty() bool {
	return q.queue.IsEmpty()
}

// Close stops publishing, published remaining elements, and tries to clean used resources.
// Please make sure this method called if this publisher need to be disposed
// or when your application exit.
func (q *Publisher) Close() error {
	q.keepPublishing = false
	q.wait.Wait()
	return nil
}

// Publishes elements periodically.
func (q *Publisher) startPublishing() {
	q.keepPublishing = true
	for {
		q.queue.lock.Lock()
		if len(q.queue.list) > 0 {
			q.publishFunc(q, q.queue.unsafePop())
		}
		q.queue.lock.Unlock()

		if !q.keepPublishing {
			break
		}

		if !q.publishingStarted {
			q.publishingStarted = true
			q.waitStart.Done()
		}

		time.Sleep(q.sleepInterval)
	}

	for _, v := range q.queue.list {
		q.publishFunc(q, v)
	}
	q.wait.Done()
}

// NewPublisher creates new publisher given a queue.
func NewPublisher(args PublisherArgs) (*Publisher, error) {
	if args.PublishFunc == nil {
		return nil, ErrPublishFunctionUnset
	}

	queue := args.Queue
	if queue == nil {
		queue = New(nil)
	}

	sleepInterval := args.SleepInterval
	if sleepInterval == 0 {
		sleepInterval = time.Millisecond * 100
	}

	p := &Publisher{
		queue:         queue,
		sleepInterval: sleepInterval,
		publishFunc:   args.PublishFunc,
	}
	p.wait.Add(1)
	p.waitStart.Add(1)
	go p.startPublishing()
	p.waitStart.Wait()

	return p, nil
}
