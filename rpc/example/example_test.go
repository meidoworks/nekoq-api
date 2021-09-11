package example_test

import (
	"testing"

	"github.com/meidoworks/nekoq-api/rpc/example"
)

func TestEmptyClientFactory_CreateClient(t *testing.T) {
	example.ExampleRpcClientUsage()
}
