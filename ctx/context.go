package ctx

type Context struct {
	RequestId    string            `ctxfield:"requestId"`
	RpcId        string            `ctxfield:"rpcId"`
	Data         map[string]string `ctxfield:"data"`
	TimeLimit    TimeLimit         `ctxfield:"timeLimit"`
	TTL          byte              `ctxfield:"ttl"`
	RequestFlags map[string]string `ctxfield:"requestFlags"`
}

type RpcContext struct {
	RequestId    string            `ctxfield:"requestId"`
	RpcId        string            `ctxfield:"rpcId"`
	Data         map[string]string `ctxfield:"data"`
	AppName      string            `ctxfield:"appName"`
	CallerNodeId string            `ctxfield:"callerNodeId"` // automatic filled by Client not by caller
	TTL          byte              `ctxfield:"ttl"`
	RouteRule    []byte            `ctxfield:"routeRule"`
}
