package echo

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/johnhaha/hakit/hadata"
)

type SubCtx struct {
	Data string
}

func (c *SubCtx) Parser(data interface{}) error {
	err := json.Unmarshal([]byte(c.Data), data)
	return err
}

//pub sub element
type pubSub struct {
	Pools map[string]chan string
	Rmt   sync.RWMutex
}

func (pb *pubSub) PubJson(data interface{}) {
	val, err := json.Marshal(data)
	if err != nil {
		log.Panic(errInvalidJson)
	}
	pb.Pub(string(val))
}

//publish data
func (pb *pubSub) Pub(data string) {
	pb.Rmt.Lock()
	defer pb.Rmt.Unlock()
	if pb.Pools == nil {
		log.Println(noSubscriber)
		return
	}
	for _, pool := range pb.Pools {
		if cap(pool) > len(pool) {
			pool <- data
		}
	}
}

//register subscriber with id and sub
func (pb *pubSub) Sub(ctx context.Context, consumer func(*SubCtx)) {
	pool := make(chan string, 2)
	id := hadata.GetStringFromInt(int(time.Now().Unix()))
	if pb.Pools == nil {
		pb.Pools = map[string]chan string{id: pool}
	} else {
		go func() {
			pb.Rmt.Lock()
			defer pb.Rmt.Unlock()
			pb.Pools[id] = pool
		}()
	}
	defer func() {
		pb.Rmt.Lock()
		defer pb.Rmt.Unlock()
		close(pool)
		delete(pb.Pools, id)
	}()
	for {
		select {
		case data := <-pool:
			ctx := SubCtx{Data: data}
			go consumer(&ctx)
		case <-ctx.Done():
			return
		}
	}
}
