package echo

import (
	"context"
	"log"
	"sync"
)

//pubsub multi channel instance for echo to work on
var channelPubSub = make(map[string]*pubSub)
var channelMt sync.RWMutex
var buffer = 6
var jobCount int = 0

//pub string data to channel
func Pub(channel string, val string) error {
	if pub, ok := channelPubSub[channel]; ok {
		return pub.Pub(val)
	}
	return errNoSubscriber
}

//pub bool data to channel
func PubBool(channel string, val bool) error {
	if pub, ok := channelPubSub[channel]; ok {
		return pub.PubBool(val)
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
	jobCount++
	log.Println("💨 sub ok, current job count is", jobCount)
	if _, ok := channelPubSub[channel]; !ok {
		channelPubSub[channel] = &pubSub{}
	}
	channelMt.Unlock()
	//set buffer to 6
	channelPubSub[channel].Sub(ctx, consumer, buffer, jobCount)
	defer func() {
		channelMt.Lock()
		defer channelMt.Unlock()
		jobCount--
		log.Println("💨 sub finished, current job count is", jobCount)
	}()
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

func SetBuffer(count int) {
	buffer = count
}
