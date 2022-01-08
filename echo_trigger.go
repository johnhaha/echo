package echo

import "context"

var triggerBuffer = 10

var echoTrigger = NewTrigger(triggerBuffer)

//🔥 this should be set before any trigger usage
func SetTriggerBuffer(buffer int) {
	if len(echoTrigger.JobRouter.Handlers) > 0 {
		panic("can not init after set router")
	}
	triggerBuffer = buffer
	echoTrigger = NewTrigger(buffer)
}

func SetTrigger(key string, handler JobHandler) {
	echoTrigger.Register(key, handler)
}

func FireTrigger(key string, data string) {
	echoTrigger.Fire(ChannelData{
		Value:   *newValue().SetValue(data),
		Channel: key,
	})
}

func TriggerStart(ctx context.Context) {
	echoTrigger.Listen(ctx)
}
