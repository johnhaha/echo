package echo

import (
	"encoding/json"
	"time"
)

type Value struct {
	ID   string
	Data string
	//created at timestamp
	Ts int64
}

func NewValue() *Value {
	return &Value{Ts: time.Now().Unix()}
}

func (v *Value) GetData() string {
	return v.Data
}

func (v *Value) GetBool() bool {
	return v.Data == BoolTrue
}

func (v *Value) GetJsonData(data any) error {
	err := json.Unmarshal([]byte(v.Data), data)
	return err
}

func (v *Value) SetValue(data string) *Value {
	v.Data = data
	return v
}

func (v *Value) SetBool(data bool) *Value {
	s := BoolFalse
	if data {
		s = BoolTrue
	}
	v.Data = s
	return v
}

func (v *Value) SetJson(data any) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.SetValue(string(d))
	return nil
}

func (v *Value) SetTime(ts int64) {
	v.Ts = ts
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

type GChannelData[T any] struct {
	Value   T
	Channel string
}

func NewGChannelData[T any](channel string, data T) GChannelData[T] {
	return GChannelData[T]{Value: data, Channel: channel}
}
