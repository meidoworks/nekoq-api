package entrance_test

import (
	"testing"

	"github.com/meidoworks/nekoq-api/entrance"
)

func TestGenerateRequestId(t *testing.T) {
	t.Log(entrance.GenerateRequestId())
}
