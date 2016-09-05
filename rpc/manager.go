package rpc

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"unsafe"

	"import.moetang.info/go/nekoq-api/errorutil"
)

var (
	clientImplMap    = make(map[string]Client)
	clientLock       = &sync.Mutex{}
	tmpClientImplMap map[string]Client

	clientFactoryMap    = make(map[string]ClientFactory)
	clientFactoryLock   = &sync.Mutex{}
	tmpClientFactoryMap map[string]ClientFactory

	initLock = &sync.Mutex{}
)

func RegisterClientFactory(name string, clientFactory ClientFactory) error {
	_, ok := clientFactoryMap[name]
	if ok {
		return errorutil.NewWithErrorCode("NEKO_RPC_DUPLICATED_CLIENT", fmt.Sprint("client:", name, "exists."))
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
	return nil
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
	for k, _ := range enabledService {
		_, exists := clientImplMap[k]
		if exists {
			continue
		}
		factory, exists := clientFactoryMap[k]
		if !exists {
			fmt.Fprintln(os.Stderr, errorutil.New("no such service mapped to: "+k))
			os.Exit(-114)
		}
		c, err := factory.CreateClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, errorutil.NewNested("create client error.", err))
			os.Exit(-115)
		}
		registerClientImpl(k, c)
	}
	initLock.Unlock()
}
