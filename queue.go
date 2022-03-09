package echo

import (
	"context"
)

type ChannelQueue struct {
	// JobRouter
	DataChannel chan ChannelData
}

//new channel queue
func NewChannelQueue(buffer int) *ChannelQueue {
	return &ChannelQueue{
		DataChannel: make(chan ChannelData, buffer),
	}
}

//insert data
func (queue *ChannelQueue) Pub(channelData ChannelData) {
	queue.DataChannel <- channelData
}

// consume
func (queue *ChannelQueue) Consume(ctx context.Context, router *JobRouter) {
	go func() {
		for data := range queue.DataChannel {
			go router.Handle(data.Channel, &SubCtx{Value: data.Value})
		}
	}()
	<-ctx.Done()
}
