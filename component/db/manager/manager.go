package manager

import (
	"sync"

	"import.moetang.info/go/nekoq-api/component/db"
	"import.moetang.info/go/nekoq-api/errorutil"
)

type DbApi interface {
	GetSimpleDb() (db.SimpleDB, error)
	GetAtomicDb() (db.AtomicDB, error)
	CloseDbApi() error
}

type DriverFactory interface {
	GetName() string
	GetDbApi(config map[string]string) (DbApi, error)
}

var (
	drivers    = make(map[string]DriverFactory)
	driverLock = new(sync.RWMutex)
)

func RegisterDriver(driverFactory DriverFactory) error {
	name := driverFactory.GetName()
	driverLock.Lock()
	defer driverLock.Unlock()
	_, ok := drivers[name]
	if ok {
		return errorutil.New("driver: " + name + " exists -> component -> nekoq-api")
	}
	drivers[name] = driverFactory
	return nil
}

func GetDbApi(name string, config map[string]string) (DbApi, error) {
	driverLock.RLock()
	defer driverLock.RUnlock()
	df, ok := drivers[name]
	if ok {
		return df.GetDbApi(config)
	} else {
		return nil, errorutil.New("driver: " + name + " not found -> component -> nekoq-api")
	}
}
