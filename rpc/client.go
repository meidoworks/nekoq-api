package rpc

import (
	"import.moetang.info/go/nekoq-api/context"
)

type Client interface {
	Call(context *context.RpcContext, request interface{}, response interface{}) error
	CallWithReturning(context *context.RpcContext, request interface{}) (interface{}, error)
}

func NewClient(name string) Client {
	//TODO
	return nil
}
