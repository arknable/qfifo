# qfifo

A simple Go FIFO queue.

## How to Use

## Queue
`Queue` is a FIFO queue implementation, to create it simply call `New()`.
```
q := qfifo.New(nil)
q.Push(1)
q.Push(2)
....
```

Under the hood, `Queue` uses slice with default capacity of 10, to increase the capacity use `QueueOptions`.
```
q = qfifo.New(&QueueOptions{
		InitialSize: 20,
	})
q.Push(1)
q.Push(2)
....
```

To add an element, use `Push()` and to pull an element out use `Pop()`. All method has internal lock so they are safe to be used in multiple go routine.

## Publisher
`Publisher` uses `Queue` and periodically pull element to be passed into a publish function. Common scenario to use `Publisher` such as logging, we want to buffer the logs and periodically write buffered logs to persistence media such as text file or database. In this case, we can do that inside the publish function and let logging methods to push logs to queue while publish function writes the logs into text file.

To create a publisher:
```
p, err = qfifo.NewPublisher(PublisherArgs{
		PublishFunc: func(p *Publisher, v interface{}) {
            // v is the pulled element
            fmt.Println(v.(int)) // assuming the queued elements are integers
        },
	})
if err != nil {
    return err
}
defer p.Close()  // make sure to call this

p.Push(1)
p.Push(2)
..........
```

`PublishFunc` is non-blocking because will be called by internal go routine. If it is unset then `NewPublisher()` returns `ErrPublishFunctionUnset`.