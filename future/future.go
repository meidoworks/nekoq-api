package future

import "time"

type Future interface {
	Wait(d time.Duration) error
	Get() interface{}
	WaitAndGet(d time.Duration) (interface{}, error)
}
