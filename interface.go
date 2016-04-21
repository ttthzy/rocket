package rocket

type RocketIFC interface {
	LoginRocketChat(uname, upass string)
	RocketChat(pd PushData)
	GetPushData(m map[string]interface{})
}
