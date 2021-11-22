package echo

import (
	"context"
	"sync"
)

//pubsub multi channel instance for echo to work on
var channelPubSub = make(map[string]*pubSub)
var channelMt sync.RWMutex

//pub string data to channel
func Pub(channel string, val string) error {
	if pub, ok := channelPubSub[channel]; ok {
		return pub.Pub(val)
	}
	return errNoSubscriber
}

//pub json encoded data
func PubJson(channel string, val interface{}) error {
	if pub, ok := channelPubSub[channel]; ok {
		return pub.PubJson(val)
	}
	return errNoSubscriber
}

//sub to some channel and take action
func Sub(ctx context.Context, channel string, consumer func(*SubCtx)) {
	channelMt.Lock()
	if _, ok := channelPubSub[channel]; !ok {
		channelPubSub[channel] = &pubSub{}
	}
	channelMt.Unlock()
	channelPubSub[channel].Sub(ctx, consumer)
}

type Suber struct {
	SubMap map[string]func(*SubCtx)
}

func NewSuber() *Suber {
	return &Suber{}
}

func (s *Suber) Add(channel string, consumer func(*SubCtx)) {
	if s.SubMap == nil {
		s.SubMap = make(map[string]func(*SubCtx))
	}
	s.SubMap[channel] = consumer

}

func (s *Suber) Sub(ctx context.Context) {
	for k, v := range s.SubMap {
		go Sub(ctx, k, v)
	}
	<-ctx.Done()
}
