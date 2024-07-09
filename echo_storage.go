package echo

import (
	"encoding/json"
	"fmt"
)

var storage Storage

func SetStorage(s Storage) {
	storage = s
}

func getTimerStorageChannel() string {
	return fmt.Sprintf("%v_%v", timerEventStorageChannel, storage.AppName())
}

func LoadTimerEvent() error {
	res, err := storage.Find(getTimerStorageChannel())
	if err != nil {
		return err
	}
	if j := len(res); j > 0 {
		event := make([]TimerEvent, j)
		for i := 0; i < j; i++ {
			data, err := storage.Get(res[i])
			if err != nil {
				return err
			}
			e := new(TimerEvent)
			err = json.Unmarshal([]byte(data), e)
			if err != nil {
				return err
			}
			//loop event run only ones if stored
			if e.Loop > 0 {
				e.Loop = 0
			}
			event[i] = *e
		}
		AddManyTimerEvent(event)
	}
	return nil
}

func storeTimerEvent(event *TimerEvent) error {
	//will not store loop event
	if event.Loop > 0 {
		return nil
	}
	if storage == nil {
		return nil
	}
	if event.ID == "" {
		id := storage.UUID()
		event.ID = id
	}
	res, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = storage.Save(event.ID, string(res))
	if err != nil {
		return err
	}
	return storage.ListAppend(getTimerStorageChannel(), event.ID)
}

func storeManyTimerEvent(event []TimerEvent) error {
	for i := 0; i < len(event); i++ {
		err := storeTimerEvent(&event[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func remTimerEvent(id string) error {
	if storage == nil {
		return nil
	}
	err := storage.Rem(id)
	if err != nil {
		return err
	}
	return storage.ListRem(getTimerStorageChannel(), id)
}
