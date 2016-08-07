package context

type RpcContext struct {
	TraceId string
	RpcId   string
	Data    map[string]string
}
