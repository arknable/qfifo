package qfifo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPublisherInstantiation(t *testing.T) {
	p, err := NewPublisher(PublisherArgs{})
	assert.Nil(t, p)
	assert.Equal(t, ErrPublishFunctionUnset, err)

	p, err = NewPublisher(PublisherArgs{
		PublishFunc: func(p *Publisher, i interface{}) {},
	})
	assert.NotNil(t, p)
	p.Close()

	assert.Equal(t, nil, err)
	assert.Equal(t, time.Millisecond*100, p.sleepInterval)
	assert.Equal(t, defaultQueueSize, cap(p.queue.list))
}

func TestPublisherPublishing(t *testing.T) {
	published := make(map[int]bool, 0)
	p, err := NewPublisher(PublisherArgs{
		PublishFunc: func(p *Publisher, v interface{}) {
			published[v.(int)] = true
		},
	})
	assert.NotNil(t, p)
	defer p.Close()

	assert.Equal(t, nil, err)

	for i := 1; i <= 10; i++ {
		published[i] = false
		p.Push(i)
	}

	p.Close()

	for i := 1; i <= 10; i++ {
		assert.True(t, published[i])
	}
}
