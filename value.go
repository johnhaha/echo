package echo

import (
	"encoding/json"
)

type Value struct {
	Data string
}

func newValue() *Value {
	return &Value{}
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

func (v *Value) SetJson(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.SetValue(string(d))
	return nil
}
