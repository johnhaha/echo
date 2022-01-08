package echo

import "context"

var queueBuffer = 10

var echoChannelQueue = NewChannelQueue(queueBuffer)

//🔥 this should be set before any queue usage
func SetQueueBuffer(b int) {
	if len(echoChannelQueue.JobRouter.Handlers) > 0 {
		panic("can not init after set router")
	}
	queueBuffer = b
	echoChannelQueue = NewChannelQueue(b)
}

//set queue handler
func SetQueueHandler(channel string, handler JobHandler) {
	echoChannelQueue.Set(channel, handler)
}

//pub string data to queue
func PubQueue(channel string, data string) {
	echoChannelQueue.Pub(ChannelData{
		Value: Value{
			Data: data,
		},
		Channel: channel,
	})
}

//pub json data to queue
func PubQueueJson(channel string, data interface{}) error {
	v := newValue()
	err := v.SetJson(data)
	if err != nil {
		return err
	}
	echoChannelQueue.Pub(ChannelData{Channel: channel, Value: *v})
	return nil
}

func ConsumeQueue(ctx context.Context) {
	echoChannelQueue.Consume(ctx)
}
