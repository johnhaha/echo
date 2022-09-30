package echo

import (
	"context"
	"sync"
	"time"
)

var echoTimerHeap = NewTimerHeap()

var heapMtx sync.RWMutex

// echo just one router
func addTimerEvent(event TimerEvent) error {
	heapMtx.Lock()
	defer heapMtx.Unlock()
	err := storeTimerEvent(&event)
	if err != nil {
		return err
	}
	echoTimerHeap.Insert(event)
	return nil
}

// add many event to timer heap, can be used in initializing heap data
func AddManyTimerEvent(event []TimerEvent) error {
	if len(event) == 0 {
		return nil
	}
	heapMtx.Lock()
	defer heapMtx.Unlock()
	err := storeManyTimerEvent(event)
	if err != nil {
		return err
	}
	echoTimerHeap.LoadMoreEvent(event)
	return nil
}

// add channel and data to timer heap
func AddTimerEvent(channel string, data string, time time.Time) {
	v := NewValue().SetValue(data)
	event := TimerEvent{
		Value:     *v,
		EventType: channel,
		Ts:        time.Unix(),
	}
	addTimerEvent(event)
}

// add loop timer event, loop in second
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

// add json data to timer heap
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

// run timer heap, this will block, sleep step is 5 second, max sleep (in second) should be bigger than 5 second
func StartTimerEventListener(ctx context.Context, maxSleep int) {
	var sleeper = NewSleeper(time.Second*5, time.Second*time.Duration(maxSleep))
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
				sleep := recent - now
				if sleep > int64(maxSleep) {
					sleep = int64(maxSleep)
				}
				time.Sleep(time.Second * time.Duration(sleep))
				continue
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
			if x.ID != "" {
				remTimerEvent(x.ID)
			}
			sleeper.Reset()
		}
	}()
	<-ctx.Done()
}
