package ApiRPC

import (
	"context"
	"time"
)

// Return a context instance with deadline
func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}
