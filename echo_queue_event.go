package echo

import "context"

var queueBuffer = 10

var echoChannelQueue = NewChannelQueue(queueBuffer)

//🔥 this should be set before any queue usage
func SetEventBuffer(b int) {
	if len(echoRouter.Handlers) > 0 {
		panic("can not init after set router")
	}
	queueBuffer = b
	echoChannelQueue = NewChannelQueue(b)
}

// echo just one router
//set queue handler
// func SetQueueHandler(channel string, handler JobHandler) {
// 	echoChannelQueue.Set(channel, handler)
// }

//pub string data to queue
func PubEvent(channel string, data string) {
	echoChannelQueue.Pub(ChannelData{
		Value: Value{
			Data: data,
		},
		Channel: channel,
	})
}

//pub json data to queue
func PubEventJson(channel string, data interface{}) error {
	v := NewValue()
	err := v.SetJson(data)
	if err != nil {
		return err
	}
	echoChannelQueue.Pub(ChannelData{Channel: channel, Value: *v})
	return nil
}

func StartEventListener(ctx context.Context) {
	echoChannelQueue.Consume(ctx, &echoRouter)
}
