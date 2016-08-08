package rpc

import (
	"import.moetang.info/go/nekoq-api/context"
	"import.moetang.info/go/nekoq-api/errors"
)

type Param struct {
	Method  string
	Request interface{}
	Context *context.RpcContext
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
	return nil, errors.New("nekoq-api: no rpc client mapped of " + name)
}
