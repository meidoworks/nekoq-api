package example

import (
	"fmt"
	"reflect"

	"goimport.moetang.info/nekoq-api/errorutil"
	"goimport.moetang.info/nekoq-api/future"
	"goimport.moetang.info/nekoq-api/rpc" // load rpc package
)

func ExampleRpcClientUsage() {
	// load implements using imports
	// first: load client factory
	rpc.RegisterClientFactory("demoService", &EmptyClientFactory{
		allowed: make(map[string]bool),
	})
	// second: load method factory
	out := ""
	rpc.RegisterMethodFactory("demoService", "echo", reflect.TypeOf(""), reflect.TypeOf(&out)) //NOTE: out must be a ptr

	// pre initiation
	rpc.InitClient()

	// usage
	c, err := rpc.GetClient("demoService")
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
	allowed map[string]bool
}

func (this *EmptyClientFactory) PreRegisterMethod(methodName string, in reflect.Type, out reflect.Type) error {
	fmt.Println("register method.", methodName)
	this.allowed[methodName] = true
	return nil
}

func (this *EmptyClientFactory) CreateClient(config map[string]string) (rpc.FullClient, error) {
	fmt.Println("new client")
	return &EmptyClient{
		allowed: this.allowed,
	}, nil
}

type EmptyClient struct {
	allowed map[string]bool
}

func (this *EmptyClient) Call(param rpc.Param, resultPtr interface{}) (timeout bool, err error) {
	_, ok := this.allowed[param.Method]
	fmt.Println(this.allowed)
	if !ok {
		return false, errorutil.New("unknown method:" + param.Method)
	}
	fmt.Println("invoke rpc:", param.Method)
	return false, nil
}

func (this *EmptyClient) AsyncCall(param rpc.Param, resultPtr interface{}) (future future.Future, err error) {
	_, ok := this.allowed[param.Method]
	if !ok {
		return nil, errorutil.New("unknown method:" + param.Method)
	}
	fmt.Println("async invoke rpc:", param.Method)
	return nil, nil
}
