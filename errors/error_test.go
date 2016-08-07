package errors_test

import (
	"testing"

	"import.moetang.info/go/nekoq-api/errors"
)

func TestError(t *testing.T) {
	e := errors.New("orig")
	e2 := errors.NewNested("outer", e)
	e3 := errors.NewNested("final", e2)

	t.Log(e)
	t.Log(e2)
	t.Log(e3)
}
