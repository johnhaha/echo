package echo

import (
	"container/heap"
	"errors"
)

type TimerEventHandler func(TimerEvent) error

type TimerEvent struct {
	Value
	EventType string
	Ts        int64
	Loop      int64
}

type TimerEventHeap struct {
	Event  []TimerEvent
	Recent int64
	Remain int
	// JobRouter
}

func NewTimerHeap() *TimerEventHeap {
	return &TimerEventHeap{}
}

func (h *TimerEventHeap) InitEvent(event []TimerEvent) {
	h.Event = event
	heap.Init(h)
	h.Update()
}

func (h *TimerEventHeap) LoadMoreEvent(event []TimerEvent) {
	h.Event = append(h.Event, event...)
	heap.Init(h)
	h.Update()
}

func (h *TimerEventHeap) Insert(event TimerEvent) {
	heap.Push(h, event)
	h.Update()
}

func (h *TimerEventHeap) Extract() (TimerEvent, error) {
	x := heap.Pop(h)
	h.Update()
	if x == nil {
		return TimerEvent{}, errors.New("found none")
	}
	event := x.(TimerEvent)
	if event.Loop > 0 {
		event.Ts += event.Loop
		h.Insert(event)
	}
	return x.(TimerEvent), nil
}

func (h *TimerEventHeap) Update() {
	l := len(h.Event)
	if l == 0 {
		h.Recent = 0
		h.Remain = 0
		return
	}
	h.Recent = h.Event[0].Ts
	h.Remain = l
}

// apply to heap interface 👇

func (h TimerEventHeap) Len() int {

	return len(h.Event)
}

func (h TimerEventHeap) Less(i, j int) bool {
	if len(h.Event) == 0 {
		return false
	}
	return h.Event[i].Ts < h.Event[j].Ts
}

func (h TimerEventHeap) Swap(i, j int) {
	if len(h.Event) == 0 {
		return
	}
	h.Event[i], h.Event[j] = h.Event[j], h.Event[i]
}

func (h *TimerEventHeap) Push(x interface{}) {
	h.Event = append(h.Event, x.(TimerEvent))
}

func (h *TimerEventHeap) Pop() interface{} {
	old := h.Event
	n := len(h.Event)
	x := old[n-1]
	h.Event = old[:n-1]
	return x
}

// apply to heap interface 👆
