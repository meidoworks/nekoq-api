package workpool_test

import (
	"fmt"
	"testing"
	"time"

	"import.moetang.info/go/nekoq-api/workpool"
)

func TestNewOrGetWorkPool(t *testing.T) {
	w := workpool.NewOrGetWorkPool("helloworld", 1)
	w.Run("demo", func() {
		fmt.Println(1)
		panic("ddd")
	})

	time.Sleep(5 * time.Second)
}
