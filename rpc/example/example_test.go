package example_test

import (
	"testing"

	"goimport.moetang.info/nekoq-api/rpc/example"
)

func TestEmptyClientFactory_CreateClient(t *testing.T) {
	example.ExampleRpcClientUsage()
}
