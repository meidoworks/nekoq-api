package tcpexample_test

import (
	"testing"
	"time"

	"goimport.moetang.info/nekoq-api/network/tcpexample"
)

func TestClientExample(t *testing.T) {
	err := tcpexample.ClientExample("localhost:6001", 10*time.Second)
	t.Log(err)
	time.Sleep(10 * time.Second)
}
