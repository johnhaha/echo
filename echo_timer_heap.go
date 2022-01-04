package echo

import (
	"sync"
	"time"
)

var echoTimerHeap = NewTimerHeap()

var heapMtx sync.RWMutex

func addEventToTimerHeap(event TimerEvent) {
	heapMtx.Lock()
	defer heapMtx.Unlock()
	echoTimerHeap.Insert(event)
}

//add many event to timer heap, can be used in initializing heap data
func AddManyEventToTimerHeap(event []TimerEvent) {
	heapMtx.Lock()
	defer heapMtx.Unlock()
	echoTimerHeap.LoadMoreEvent(event)
}

//add channel and data to timer heap
func AddToTimerHeap(channel string, data string, time time.Time) {
	event := TimerEvent{
		Value:     *newValue(data),
		EventType: channel,
		Ts:        time.Unix(),
	}
	addEventToTimerHeap(event)
}

//add json data to timer heap
func AddJsonDataToTimerHeap(channel string, data interface{}, time time.Time) error {
	v := newValue("")
	err := v.SetJson(data)
	if err != nil {
		return err
	}
	event := TimerEvent{
		Value:     *v,
		EventType: channel,
		Ts:        time.Unix(),
	}
	addEventToTimerHeap(event)
	return nil
}

//set timer heap handler
func SetTimerHeapHandler(channel string, handler JobHandler) {
	echoTimerHeap.Set(channel, handler)
}

//run timer heap, this will block
func RunTimerHeap() {
	sleeper := NewSleeper(time.Second*5, time.Second*300)
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
		go echoTimerHeap.Handle(x.EventType, &SubCtx{
			Value: x.Value,
		})
		sleeper.Reset()
	}
}
