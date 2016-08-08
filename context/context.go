package context

type Context struct {
	TraceId string
	RpcId   string
	Data    map[string]string
	AppName string
}

type RpcContext struct {
	TraceId string
	RpcId   string
	Data    map[string]string
	AppName string
}
