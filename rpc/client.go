package rpc

import (
	"import.moetang.info/go/nekoq-api/ctx"
	"import.moetang.info/go/nekoq-api/errorutil"
	"import.moetang.info/go/nekoq-api/future"
)

type Param struct {
	Method  string
	Request interface{}
	Context *ctx.RpcContext
}

type Result struct {
	Response     interface{}
	CalleeNodeId string
}

type rpcCall interface {
	Call(param Param, response interface{}) (timeout bool, err error)
	CallReturning(param Param) (result interface{}, timeout bool, err error)
	AsyncCallReturning(param Param) (future future.Future, timeout bool, err error)
}

type Client interface {
	rpcCall
}

type ClientFactory interface {
	CreateClient() (Client, error)
}

func NewClient(name string) (Client, error) {
	c, ok := clientImplMap[name]
	if ok {
		return c, nil
	}
	return nil, errorutil.New("nekoq-api: no rpc client mapped of " + name)
}
