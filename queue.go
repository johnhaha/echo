package echo

import (
	"context"
	"sync"
)

type ChannelQueue struct {
	JobRouter
	DataChannel chan ChannelData
	Mtx         sync.RWMutex
}

//new channel queue
func NewChannelQueue(buffer int) *ChannelQueue {
	return &ChannelQueue{
		JobRouter:   JobRouter{Handlers: make(map[string]JobHandler)},
		DataChannel: make(chan ChannelData, buffer),
	}
}

//insert data
func (queue *ChannelQueue) Pub(channelData ChannelData) {
	queue.Mtx.Lock()
	defer queue.Mtx.Unlock()
	queue.DataChannel <- channelData
}

// consume
func (queue *ChannelQueue) Consume(ctx context.Context) {
	go func() {
		for data := range queue.DataChannel {
			go queue.JobRouter.Handle(data.Channel, &SubCtx{Value: data.Value})
		}
	}()
	<-ctx.Done()
}
