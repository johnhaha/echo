package echo

//used in queue and timer heap
var echoRouter JobRouter

func SetEventHandler(channel string, handler JobHandler) {
	echoRouter.Set(channel, handler)
}
