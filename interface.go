package rocket

type RocketIFC interface {
	getRocketUserToken()
	PushRocketChat()
	GetPushData(m map[string]interface{})
}
