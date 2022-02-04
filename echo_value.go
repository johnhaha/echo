package echo

import (
	"sync"
)

var valueMap = make(map[string]*Value)
var valueMtx sync.RWMutex

func SetValue(k string, v string) {
	value := NewValue().SetValue(v)
	valueMtx.Lock()
	valueMap[k] = value
	valueMtx.Unlock()
}

func SetBoolValue(k string, v bool) {
	value := NewValue()
	value.SetBool(v)
	valueMtx.Lock()
	valueMap[k] = value
	valueMtx.Unlock()
}

func SetJsonValue(k string, data any) error {
	value := NewValue()
	err := value.SetJson(data)
	if err != nil {
		return err
	}
	valueMtx.Lock()
	valueMap[k] = value
	valueMtx.Unlock()
	return nil
}

func GetValue(k string) (string, error) {
	valueMtx.RLock()
	defer valueMtx.RUnlock()
	if v, ok := valueMap[k]; ok {
		return v.GetData(), nil
	}
	return "", errNotFound
}

func GetBoolValue(k string) (bool, error) {
	valueMtx.RLock()
	defer valueMtx.RUnlock()
	if v, ok := valueMap[k]; ok {
		return v.GetBool(), nil
	}
	return false, errNotFound
}

func GetJsonValue(k string, data any) error {
	valueMtx.RLock()
	defer valueMtx.RUnlock()
	if v, ok := valueMap[k]; ok {
		return v.GetJsonData(data)
	}
	return errNotFound
}

func RemValue(k string) {
	valueMtx.RLock()
	defer valueMtx.RUnlock()
	delete(valueMap, k)
}
