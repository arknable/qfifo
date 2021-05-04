package qfifo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueInstantiation(t *testing.T) {
	q := New(nil)
	assert.NotNil(t, q)
	assert.Equal(t, 0, len(q.list))
	assert.Equal(t, defaultSize, cap(q.list))

	q = New(&QueueOptions{
		InitialSize: 20,
	})
	assert.NotNil(t, q)
	assert.Equal(t, 0, len(q.list))
	assert.Equal(t, 20, cap(q.list))
}

func TestQueuePushAndPop(t *testing.T) {
	q := New(nil)
	assert.NotNil(t, q)

	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	assert.Equal(t, 5, len(q.list))
	assert.Equal(t, defaultSize, cap(q.list))
	for i := 0; i < 5; i++ {
		assert.Equal(t, i+1, q.list[i])
	}

	refMap := make(map[int]bool)
	for i := 1; i <= 5; i++ {
		refMap[i] = false
	}

	var v int
	for {
		if q.IsEmpty() {
			break
		}

		v = q.Pop().(int)
		refMap[v] = true
	}

	isVerified := true
	for _, state := range refMap {
		isVerified = isVerified && state
	}
	assert.True(t, isVerified)
}

func TestQueueClear(t *testing.T) {
	q := New(nil)
	assert.NotNil(t, q)

	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	assert.Equal(t, 5, len(q.list))
	assert.Equal(t, defaultSize, cap(q.list))

	q.Clear()
	assert.Equal(t, 0, len(q.list))
	assert.Equal(t, defaultSize, cap(q.list))
	assert.True(t, q.IsEmpty())
}
