package rpc

import (
	"sync"

	"github.com/meidoworks/nekoq-api/errorutil"
)

var serverFactories = make(map[string]ServerFactory)
var mLock = &sync.Mutex{}

type ServerFactory interface {
	Start() error
	Stop() error
}

func RegisterServerFactory(serviceName string, serverFactory ServerFactory) error {
	mLock.Lock()
	_, ok := serverFactories[serviceName]
	if ok {
		mLock.Unlock()
		return errorutil.New("same server name has registered: " + serviceName)
	}
	serverFactories[serviceName] = serverFactory
	mLock.Unlock()
	return nil
}

func StartAllServices() error {
	for k, v := range serverFactories {
		err := v.Start()
		if err != nil {
			return errorutil.NewNested("start service error: "+k, err)
		}
	}
	return nil
}
