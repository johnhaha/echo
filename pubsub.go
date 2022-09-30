package echo

import (
	"context"
	"sync"
)

// pub sub stateless single block
type PubSub[T any] struct {
	//use map chan to support multiple subscription
	Pools map[string]chan T
	Rmt   sync.RWMutex
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{}
}

// publish data
func (pb *PubSub[T]) Pub(data T) error {
	pb.Rmt.Lock()
	defer pb.Rmt.Unlock()
	if pb.Pools == nil {
		return errNoSubscriber
	}
	for _, pool := range pb.Pools {
		if cap(pool) > len(pool) {
			pool <- data
		}
	}
	return nil
}

// register subscriber with id and sub
func (pb *PubSub[T]) Sub(ctx context.Context, group string, buffer int, consumer func(T)) {
	pool := make(chan T, buffer)
	pb.Rmt.Lock()
	if pb.Pools == nil {
		pb.Pools = make(map[string]chan T)
	}
	pb.Pools[group] = pool
	pb.Rmt.Unlock()
	defer func() {
		pb.Rmt.Lock()
		defer pb.Rmt.Unlock()
		close(pool)
		delete(pb.Pools, group)
	}()
	for {
		select {
		case data := <-pool:
			go consumer(data)
		case <-ctx.Done():
			return
		}
	}
}
