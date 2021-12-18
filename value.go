package echo

import (
	"encoding/json"
)

type Value struct {
	Data string
}

func newValue(v string) *Value {
	return &Value{Data: v}
}

func (v *Value) GetData() string {
	return v.Data
}

func (v *Value) GetBool() bool {
	return v.Data == BoolTrue
}

func (v *Value) GetJsonData(data interface{}) error {
	err := json.Unmarshal([]byte(v.Data), data)
	return err
}

func (v *Value) SetValue(data string) {
	v.Data = data
}

func (v *Value) SetBool(data bool) {
	s := BoolFalse
	if data {
		s = BoolTrue
	}
	v.Data = s
}

func (v *Value) SetJson(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.SetValue(string(d))
	return nil
}
