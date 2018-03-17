package rpc

import (
	"reflect"

	"goimport.moetang.info/nekoq-api/ctx"
	"goimport.moetang.info/nekoq-api/errorutil"
	"goimport.moetang.info/nekoq-api/future"
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

type FullClient interface {
	Client
}

type ClientFactory interface {
	PreRegisterMethod(methodName string, in reflect.Type, out reflect.Type) error
	CreateClient(config map[string]string) (FullClient, error)
}

func GetClient(name string) (Client, error) {
	c, ok := clientImplMap[name]
	if ok {
		return c, nil
	}
	return nil, errorutil.New("nekoq-api: no rpc client mapped of " + name)
}
