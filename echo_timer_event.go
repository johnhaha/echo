package echo

import (
	"context"
	"sync"
	"time"
)

var echoTimerHeap = NewTimerHeap()

var heapMtx sync.RWMutex

// echo just one router
func addTimerEvent(event TimerEvent) {
	heapMtx.Lock()
	defer heapMtx.Unlock()
	echoTimerHeap.Insert(event)
}

//add many event to timer heap, can be used in initializing heap data
func AddManyTimerEvent(event []TimerEvent) {
	heapMtx.Lock()
	defer heapMtx.Unlock()
	echoTimerHeap.LoadMoreEvent(event)
}

//add channel and data to timer heap
func AddTimerEvent(channel string, data string, time time.Time) {
	v := NewValue().SetValue(data)
	event := TimerEvent{
		Value:     *v,
		EventType: channel,
		Ts:        time.Unix(),
	}
	addTimerEvent(event)
}

//add loop timer event, loop in second
func AddLoopTimerEvent(channel string, data string, time time.Time, loop int64) {
	v := NewValue().SetValue(data)
	event := TimerEvent{
		Value:     *v,
		EventType: channel,
		Ts:        time.Unix(),
		Loop:      loop,
	}
	addTimerEvent(event)
}

//add json data to timer heap
func AddJsonDataToTimerEvent(channel string, data any, time time.Time) error {
	v := NewValue()
	err := v.SetJson(data)
	if err != nil {
		return err
	}
	event := TimerEvent{
		Value:     *v,
		EventType: channel,
		Ts:        time.Unix(),
	}
	addTimerEvent(event)
	return nil
}

//run timer heap, this will block
func StartTimerEventListener(ctx context.Context) {
	sleeper := NewSleeper(time.Second*5, time.Second*300)
	go func() {
		for {
			now := time.Now().Unix()
			heapMtx.RLock()
			recent := echoTimerHeap.Recent
			heapMtx.RUnlock()
			if recent == 0 {
				sleeper.Sleep()
				continue
			}
			if recent-now > 0 {
				time.Sleep(time.Second * time.Duration(recent-now))
			}
			heapMtx.Lock()
			x, err := echoTimerHeap.Extract()
			heapMtx.Unlock()
			if err != nil {
				continue
			}
			go echoRouter.Handle(x.EventType, &SubCtx{
				Value: x.Value,
			})
			sleeper.Reset()
		}
	}()
	<-ctx.Done()
}
