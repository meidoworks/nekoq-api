package errorutil_test

import (
	"testing"

	"github.com/meidoworks/nekoq-api/errorutil"
)

func TestError(t *testing.T) {
	e := errorutil.New("orig")
	e2 := errorutil.NewNested("outer", e)
	e3 := errorutil.NewNested("final", e2)

	t.Log(e)
	t.Log(e2)
	t.Log(e3)
}
