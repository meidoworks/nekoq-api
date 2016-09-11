package rpc

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"

	"import.moetang.info/go/nekoq-api/errorutil"
)

var (
	clientImplMap    = make(map[string]Client)
	clientLock       = &sync.Mutex{}
	tmpClientImplMap map[string]Client

	methodMap = make(map[string]map[string]bool)

	clientFactoryMap    = make(map[string]ClientFactory)
	clientFactoryLock   = &sync.Mutex{}
	tmpClientFactoryMap map[string]ClientFactory

	initLock = &sync.Mutex{}
)

func RegisterMethodFactory(name string, method string, in reflect.Type, out reflect.Type) {
	clientFactoryLock.Lock()
	var v map[string]bool
	var ok bool
	if v, ok = methodMap[name]; !ok {
		v = make(map[string]bool)
		methodMap[name] = v
	}
	if _, exist := v[method]; exist {
		clientFactoryLock.Unlock()
		panic(errorutil.NewWithErrorCode("NEKO_RPC_DUPLICATED_CLIENT_METHOD", fmt.Sprint("client:", name, "method:", method, "exists.")))
	}
	if out.Kind() != reflect.Ptr {
		clientFactoryLock.Unlock()
		panic(errorutil.New("NEKO_RPC_OUT_PARAM_TYPE_MUST_BE_POINTER"))
	}
	cf, ok := clientFactoryMap[name]
	if !ok {
		clientFactoryLock.Unlock()
		panic(errorutil.NewWithErrorCode("NEKO_RPC_NO_CLIENT_FOUND", fmt.Sprint("client:", name, "not found.")))
	}
	cf.PreRegisterMethod(method, in, out)
	v[method] = true
	clientFactoryLock.Unlock()
}

func RegisterClientFactory(name string, clientFactory ClientFactory) {
	_, ok := clientFactoryMap[name]
	if ok {
		panic(errorutil.NewWithErrorCode("NEKO_RPC_DUPLICATED_CLIENT", fmt.Sprint("client:", name, "exists.")))
	}
	clientFactoryLock.Lock()
	// must be assigned to a global field
	tmpClientFactoryMap = make(map[string]ClientFactory)
	for k, v := range clientFactoryMap {
		tmpClientFactoryMap[k] = v
	}
	tmpClientFactoryMap[name] = clientFactory
	orgiPtr := (*uintptr)(unsafe.Pointer(&clientFactoryMap))
	newPtr := (*uintptr)(unsafe.Pointer(&tmpClientFactoryMap))
	atomic.StoreUintptr(orgiPtr, *newPtr)
	tmpClientFactoryMap = nil
	clientFactoryLock.Unlock()
}

func registerClientImpl(name string, client Client) {
	clientLock.Lock()
	// must be assigned to a global field
	tmpClientImplMap = make(map[string]Client)
	for k, v := range clientImplMap {
		tmpClientImplMap[k] = v
	}
	tmpClientImplMap[name] = client
	orgiPtr := (*uintptr)(unsafe.Pointer(&clientImplMap))
	newPtr := (*uintptr)(unsafe.Pointer(&tmpClientImplMap))
	atomic.StoreUintptr(orgiPtr, *newPtr)
	tmpClientImplMap = nil
	clientLock.Unlock()
}

func InitClient() {
	initLock.Lock()
	for k, v := range enabledService {
		_, exists := clientImplMap[k]
		if exists {
			continue
		}
		factory, exists := clientFactoryMap[k]
		if !exists {
			fmt.Fprintln(os.Stderr, errorutil.New("no such service mapped to: "+k))
			os.Exit(-114)
		}

		// check method registered!!!
		mMap := methodMap[k]
		for methodName, _ := range v {
			if _, ok := mMap[methodName]; !ok {
				fmt.Fprintln(os.Stderr, errorutil.New("no such method:"+methodName+" mapped to: "+k))
				os.Exit(-116)
			}
		}

		c, err := factory.CreateClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, errorutil.NewNested("create client error.", err))
			os.Exit(-115)
		}
		registerClientImpl(k, c)
	}
	log.Println("init rpc client...")
	initLock.Unlock()
}
