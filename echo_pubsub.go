package echo

import (
	"context"
	"log"
	"sync"
)

var channelPubSub = make(map[string]*pubSub)
var channelMt sync.RWMutex

//pub string data to channel
func Pub(channel string, val string) {
	if pub, ok := channelPubSub[channel]; ok {
		pub.Pub(val)
	} else {
		log.Println(noSubscriber)
	}
}

//pub json encoded data
func PubJson(channel string, val interface{}) {
	if pub, ok := channelPubSub[channel]; ok {
		pub.PubJson(val)
	} else {
		log.Println(noSubscriber)
	}
}

//sub to some channel and take action
func Sub(ctx context.Context, channel string, consumer func(string)) {
	channelMt.Lock()
	if _, ok := channelPubSub[channel]; !ok {
		channelPubSub[channel] = &pubSub{}
	}
	channelMt.Unlock()
	channelPubSub[channel].Sub(ctx, consumer)
}
