package ctx

type ReqId struct {
	RequestId    string
	CurrentRpcId string
	NextSubRpcNo int
}

type Context struct {
	ReqId     ReqId
	Data      map[string]string
	TimeLimit TimeLimit
	TTL       TTL
}

func NewContext(reqId, currentRpcId string, timelimit TimeLimit, ttl TTL) *Context {
	return &Context{
		ReqId: ReqId{
			RequestId:    reqId,
			CurrentRpcId: currentRpcId,
			NextSubRpcNo: 1,
		},
		Data:      make(map[string]string),
		TimeLimit: timelimit,
		TTL:       ttl,
	}
}

type RpcContext struct {
	RequestId    string            `ctxfield:"requestId"`
	RpcId        string            `ctxfield:"rpcId"`
	Data         map[string]string `ctxfield:"data"`
	AppName      string            `ctxfield:"appName"`
	CallerNodeId string            `ctxfield:"callerNodeId"` // automatic filled by Client not by caller
	TTL          byte              `ctxfield:"ttl"`
	MaxTTL       byte              `ctxfield:"maxTtl"`
	RouteRule    []byte            `ctxfield:"routeRule"`
}
