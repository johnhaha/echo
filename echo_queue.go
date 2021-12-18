package echo

import (
	"context"
	"sync"
)

var channelQueue = make(map[string]*Queue)

var channelQueueMtx sync.RWMutex

type Streamer struct {
	QueueHandler map[string]func(*SubCtx)
}

func NewStreamer() *Streamer {
	return &Streamer{QueueHandler: make(map[string]func(*SubCtx))}
}

func (s *Streamer) Add(name string, handler JobHandler) {
	channelQueueMtx.Lock()
	defer channelQueueMtx.Unlock()
	if _, ok := channelQueue[name]; ok {
		return
	}
	queue := newQueue()
	channelQueue[name] = queue
	if s.QueueHandler == nil {
		s.QueueHandler = make(map[string]func(*SubCtx))
	}
	s.QueueHandler[name] = handler
}

func (s *Streamer) Stream(ctx context.Context) {
	for k, v := range channelQueue {
		go v.Consume(ctx, s.QueueHandler[k])
	}
	<-ctx.Done()
}

func PubQ(channel string, data string) {
	channelQueueMtx.RLock()
	defer channelQueueMtx.RUnlock()
	if q, ok := channelQueue[channel]; ok {
		q.Append(data)
		return
	}
	queue := newQueue()
	channelQueue[channel] = queue
	queue.Append(data)
}

func PubBoolQ(channel string, data bool) {
	channelQueueMtx.RLock()
	defer channelQueueMtx.RUnlock()
	if q, ok := channelQueue[channel]; ok {
		q.AppendBool(data)
		return
	}
	queue := newQueue()
	channelQueue[channel] = queue
	queue.AppendBool(data)
}

func PubJsonQ(channel string, data interface{}) error {
	channelQueueMtx.RLock()
	defer channelQueueMtx.RUnlock()
	if q, ok := channelQueue[channel]; ok {
		q.AppendJson(data)
		return nil
	}
	queue := newQueue()
	channelQueue[channel] = queue
	queue.AppendJson(data)
	return nil
}
