package entrance

import (
	"fmt"
	"time"

	"goimport.moetang.info/nekoq-api/env"
)

func GenerateRequestId() string {
	t := time.Now()
	return fmt.Sprint(env.GetNodeId(), ".", t.Unix(), ".", t.UnixNano())
}
