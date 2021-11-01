package echo

import (
	"context"
	"sync"
)

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
