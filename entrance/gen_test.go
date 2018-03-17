package entrance_test

import (
	"testing"

	"goimport.moetang.info/nekoq-api/entrance"
)

func TestGenerateRequestId(t *testing.T) {
	t.Log(entrance.GenerateRequestId())
}
