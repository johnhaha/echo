package echo

import (
	"context"
	"log"
	"sync"
)

type Trigger struct {
	JobRouter
	Data chan ChannelData
	Mtx  sync.RWMutex
}

func NewTrigger(buffer int) *Trigger {
	return &Trigger{Data: make(chan ChannelData, buffer)}
}

func (trigger *Trigger) Register(key string, handler JobHandler) {
	trigger.Mtx.Lock()
	defer trigger.Mtx.Unlock()
	trigger.Set(key, handler)
}

func (trigger *Trigger) Fire(data ChannelData) {
	trigger.Data <- data
}

func (trigger *Trigger) Listen(ctx context.Context) {
	go func() {
		for data := range trigger.Data {
			trigger.Handle(data.Channel, &SubCtx{
				Value: data.Value,
			})
			trigger.Done(data.Channel)
		}
	}()
	<-ctx.Done()
}

func (trigger *Trigger) Done(key string) {
	trigger.Mtx.Lock()
	defer trigger.Mtx.Unlock()
	trigger.Set(key, func(sc *SubCtx) {
		log.Println("🔥 trigger has been fired")
	})
}
