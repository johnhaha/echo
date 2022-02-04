package echo

import (
	"errors"
	"sync"
)

type SubCtx struct {
	Value
}

func (c *SubCtx) Parser(data any) error {
	err := c.GetJsonData(data)
	return err
}

func GetData[T any](ctx SubCtx) (T, error) {
	var data T
	err := ctx.Parser(&data)
	return data, err
}

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
