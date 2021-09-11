package entrance

import (
	"fmt"
	"time"

	"github.com/meidoworks/nekoq-api/env"
)

func GenerateRequestId() string {
	t := time.Now()
	return fmt.Sprint(env.GetNodeId(), ".", t.Unix(), ".", t.UnixNano())
}
