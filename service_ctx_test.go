package go_service_mgr

import (
	"fmt"
	"testing"
	"time"
)

func TestServiceContext(t *testing.T) {
	_ = ServiceContext()
	fmt.Println("TestServiceContext")
}

func TestServiceContextWithValue(t *testing.T) {
	ctx := ServiceContextWithValue("key", "value")
	val := ctx.Value("key").(string)
	fmt.Println("TestServiceContextWithValue", val)
}

func TestServiceContextWithCancel1(t *testing.T) {
	ctx, cbCancel := ServiceContextWithCancel()

	go func() {
		index := 1
		for {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("loop...")
			if index >= 3 {
				fmt.Println("cancel...")
				cbCancel()
			}
			index++
		}
	}()

	select {
	case <-time.After(1000 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("halt")
	}
}

func TestServiceContextWithCancel2(t *testing.T) {
	ctx, cbCancel := ServiceContextWithCancel()

	go func() {
		index := 1
		for {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("loop...")
			if index >= 3 {
				fmt.Println("cancel...")
				cbCancel()
			}
			index++
		}
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("halt")
	}
}

func TestServiceContextWithTimeout1(t *testing.T) {
	ctx, _ := ServiceContextWithTimeout(1)

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}

func TestServiceContextWithTimeout2(t *testing.T) {
	ctx, _ := ServiceContextWithTimeout(1)

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}

func TestServiceContextWithDeadline1(t *testing.T) {
	ctx, _ := ServiceContextWithDeadline(time.Now().Add(time.Second))

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}

func TestServiceContextWithDeadline2(t *testing.T) {
	ctx, _ := ServiceContextWithDeadline(time.Now().Add(time.Second))

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}
