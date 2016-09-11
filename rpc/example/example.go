package example

import (
	"fmt"
	"reflect"

	"import.moetang.info/go/nekoq-api/future"
	"import.moetang.info/go/nekoq-api/rpc" // load rpc package
)

func ExampleRpcClientUsage() {
	// load implements using imports
	// first: load client factory
	rpc.RegisterClientFactory("demoService", &EmptyClientFactory{})
	// second: load method factory
	out := ""
	rpc.RegisterMethodFactory("demoService", "echo", reflect.TypeOf(""), reflect.TypeOf(&out)) //NOTE: out must be a ptr

	// pre initiation
	rpc.InitClient()

	// usage
	c, err := rpc.NewClient("demoService")
	if err != nil {
		fmt.Println("new client error.", err)
	}
	var v string
	timeout, err := c.Call(rpc.Param{
		Method: "echo",
	}, &v)
	fmt.Println(timeout, err)
}

type EmptyClientFactory struct {
}

func (this *EmptyClientFactory) PreRegisterMethod(methodName string, in reflect.Type, out reflect.Type) error {
	fmt.Println("register method.", methodName)
	return nil
}

func (this *EmptyClientFactory) CreateClient() (rpc.Client, error) {
	fmt.Println("new client")
	return &EmptyClient{}, nil
}

type EmptyClient struct {
}

func (this *EmptyClient) Call(param rpc.Param, resultPtr interface{}) (timeout bool, err error) {
	fmt.Println("invoke rpc:", param.Method)
	return false, nil
}

func (this *EmptyClient) AsyncCall(param rpc.Param, resultPtr interface{}) (future future.Future, err error) {
	fmt.Println("async invoke rpc:", param.Method)
	return nil, nil
}
