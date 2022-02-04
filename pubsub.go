package echo

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
)

//pub sub staleless single block
type pubSub struct {
	//use map chan to support multiple subscription
	Pools map[string]chan string
	Rmt   sync.RWMutex
}

func (pb *pubSub) PubJson(data any) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return pb.Pub(string(val))
}

//publish data
func (pb *pubSub) Pub(data string) error {
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

//publish data
func (pb *pubSub) PubBool(data bool) error {
	pb.Rmt.Lock()
	defer pb.Rmt.Unlock()
	if pb.Pools == nil {
		return errNoSubscriber
	}
	for _, pool := range pb.Pools {
		if cap(pool) > len(pool) {
			s := BoolFalse
			if data {
				s = BoolTrue
			}
			pool <- s
		}
	}
	return nil
}

//register subscriber with id and sub
func (pb *pubSub) Sub(ctx context.Context, consumer func(*SubCtx), buffer int, count int) {
	pool := make(chan string, buffer)
	pb.Rmt.Lock()
	id := strconv.Itoa(count)
	if pb.Pools == nil {
		pb.Pools = make(map[string]chan string)
	}
	pb.Pools[id] = pool
	pb.Rmt.Unlock()
	defer func() {
		pb.Rmt.Lock()
		defer pb.Rmt.Unlock()
		close(pool)
		delete(pb.Pools, id)
	}()
	for {
		select {
		case data := <-pool:
			subCtx := SubCtx{Value: *NewValue().SetValue(data)}
			go consumer(&subCtx)
		case <-ctx.Done():
			return
		}
	}
}
