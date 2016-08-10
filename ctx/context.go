package ctx

type Context struct {
	TraceId   string
	RpcId     string
	Data      map[string]string
	TimeLimit TimeLimit
}

type RpcContext struct {
	TraceId      string
	RpcId        string
	Data         map[string]string
	AppName      string
	CallerNodeId string // automatic filled by Client not by caller
}
