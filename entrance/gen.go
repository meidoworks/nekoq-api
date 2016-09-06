package entrance

import (
	"fmt"
	"time"

	"import.moetang.info/go/nekoq-api/env"
)

func GenerateRequestId() string {
	t := time.Now()
	return fmt.Sprint(env.GetNodeId(), ".", t.Unix(), ".", t.UnixNano())
}
