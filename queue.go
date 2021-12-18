package echo

import (
	"context"
	"encoding/json"
	"errors"
)

type Queue struct {
	Stream chan string
}

func newQueue() *Queue {
	return &Queue{Stream: make(chan string, 10)}
}

func (q *Queue) Append(data string) error {
	if len(q.Stream) < cap(q.Stream) {
		q.Stream <- data
		return nil
	}
	return errors.New("💧 queue is full")
}

func (q *Queue) AppendJson(data interface{}) error {
	if len(q.Stream) < cap(q.Stream) {
		d, err := json.Marshal(data)
		if err != nil {
			return err
		}
		q.Append(string(d))
		return nil
	}
	return errors.New("💧 queue is full")
}

func (q *Queue) AppendBool(data bool) error {
	if len(q.Stream) < cap(q.Stream) {
		s := BoolFalse
		if data {
			s = BoolTrue
		}
		q.Append(string(s))
		return nil
	}
	return errors.New("💧 queue is full")
}

func (q *Queue) Consume(ctx context.Context, consumer func(*SubCtx)) {
	for {
		select {
		case data := <-q.Stream:
			subCtx := &SubCtx{
				Value: *newValue(data),
			}
			go consumer(subCtx)
		case <-ctx.Done():
			return
		}

	}
}
