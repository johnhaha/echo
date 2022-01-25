package echo

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type SubCtx struct {
	Value
}

func (c *SubCtx) Parser(data interface{}) error {
	err := c.GetJsonData(data)
	return err
}

const (
	BoolTrue  = "True"
	BoolFalse = "False"
)

type JobHandler func(*SubCtx)

type JobRouter struct {
	Handlers map[string]JobHandler
	mtx      sync.RWMutex
}

func (r *JobRouter) Set(channel string, handler JobHandler) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.Handlers == nil {
		r.Handlers = make(map[string]JobHandler)
	}
	r.Handlers[channel] = handler
}

func (r *JobRouter) Handle(channel string, ctx *SubCtx) error {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if r.Handlers == nil {
		return errors.New("no handlers found")
	}
	if handler, ok := r.Handlers[channel]; ok {
		handler(ctx)
		return nil
	}
	return errors.New("no handlers found")
}

type Sleeper struct {
	Duration time.Duration
	Step     time.Duration
	Max      time.Duration
}

func NewSleeper(step time.Duration, max time.Duration) *Sleeper {
	return &Sleeper{
		Duration: step,
		Step:     step,
		Max:      max,
	}
}

func (sleeper *Sleeper) Sleep() {
	time.Sleep(sleeper.Duration)
	if sleeper.Duration < sleeper.Max {
		sleeper.Duration += sleeper.Step
	}
}

func (sleeper *Sleeper) Reset() {
	sleeper.Duration = sleeper.Step
}

type ChannelData struct {
	Value
	Channel string
}

func GetChannelDataFromJson(j string) (*ChannelData, error) {
	var channelData ChannelData
	err := json.Unmarshal([]byte(j), &channelData)
	if err != nil {
		return nil, err
	}
	return &channelData, nil
}
