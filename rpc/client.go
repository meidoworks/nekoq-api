package rpc

import (
	"reflect"

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
	Call(param Param, resultPtr interface{}) (timeout bool, err error)
	AsyncCall(param Param, resultPtr interface{}) (future future.Future, err error)
}

type Client interface {
	rpcCall
}

type ClientFactory interface {
	PreRegisterMethod(methodName string, in reflect.Type, out reflect.Type) error
	CreateClient(config map[string]string) (Client, error)
}

func NewClient(name string) (Client, error) {
	c, ok := clientImplMap[name]
	if ok {
		return c, nil
	}
	return nil, errorutil.New("nekoq-api: no rpc client mapped of " + name)
}
