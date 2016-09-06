package entrance_test

import (
	"testing"

	"import.moetang.info/go/nekoq-api/entrance"
)

func TestGenerateRequestId(t *testing.T) {
	t.Log(entrance.GenerateRequestId())
}
