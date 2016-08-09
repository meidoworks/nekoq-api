package rpc

import (
	"import.moetang.info/go/nekoq-api/ctx"
	"import.moetang.info/go/nekoq-api/errorutil"
)

type Param struct {
	Method  string
	Request interface{}
	Context *ctx.RpcContext
}

type Client interface {
	Call(param Param, response interface{}) error
	CallWithReturning(param Param) (interface{}, error)
}

func NewClient(name string) (Client, error) {
	c, ok := clientImplMap[name]
	if ok {
		return c, nil
	}
	return nil, errorutil.New("nekoq-api: no rpc client mapped of " + name)
}
