package rpc

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	clientImplMap    = make(map[string]Client)
	clientLock       = &sync.Mutex{}
	tmpClientImplMap map[string]Client
)

func RegisterClientImpl(name string, client Client) {
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
