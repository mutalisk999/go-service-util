package go_service_mgr

import (
	"fmt"
	"testing"
)

func TestServiceContext(t *testing.T) {
	_ = ServiceContext()
	fmt.Println("TestServiceContext")
}

func TestServiceContextWithCancel(t *testing.T) {

}

func TestServiceContextWithTimeout(t *testing.T) {

}

func TestServiceContextWithDeadline(t *testing.T) {

}
