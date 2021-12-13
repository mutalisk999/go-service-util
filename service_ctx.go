package go_service_mgr

import (
	"context"
	"time"
)

func ServiceContext() context.Context {
	return context.Background()
}

func ServiceContextWithValue(key interface{}, value interface{}) context.Context {
	ctx := context.WithValue(context.Background(), key, value)
	return ctx
}

func ServiceContextWithCancel() (context.Context, context.CancelFunc) {
	ctx, cbCancel := context.WithCancel(context.Background())
	return ctx, cbCancel
}

func ServiceContextWithTimeout(timeout uint64) (context.Context, context.CancelFunc) {
	ctx, cbCancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	return ctx, cbCancel
}

func ServiceContextWithDeadline(deadline time.Time) (context.Context, context.CancelFunc) {
	ctx, cbCancel := context.WithDeadline(context.Background(), deadline)
	return ctx, cbCancel
}
