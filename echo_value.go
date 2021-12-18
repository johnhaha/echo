package echo

import (
	"sync"
)

var valueMap = make(map[string]*Value)
var valueMtx sync.RWMutex

func SetValue(k string, v string) {
	value := newValue(v)
	valueMtx.Lock()
	valueMap[k] = value
	valueMtx.Unlock()
}

func SetBoolValue(k string, v bool) {
	value := newValue("")
	value.SetBool(v)
	valueMtx.Lock()
	valueMap[k] = value
	valueMtx.Unlock()
}

func SetJsonValue(k string, data interface{}) {
	value := newValue("")
	value.SetJson(data)
	valueMtx.Lock()
	valueMap[k] = value
	valueMtx.Unlock()
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

func GetJsonValue(k string, data interface{}) error {
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
