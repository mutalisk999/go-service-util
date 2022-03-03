package go_service_util

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestCreateRegSvcClient(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	svcCli, err := CreateRegSvcClient([]string{"http://127.0.0.1:12379", "http://127.0.0.1:22379", "http://127.0.0.1:32379"}, nil, nil, 3)
	if err != nil {
		logger.Println("CreateRegSvcClient, err:", err)
		return
	}
	defer svcCli.DisposeRegSvcClient()

	cbCancel, err := svcCli.RegisterService(5, 30, "", "test", "key", "value", logger)
	if err != nil {
		logger.Println("RegisterService, err:", err)
		return
	}

	time.Sleep(300 * time.Second)
	cbCancel()
	time.Sleep(10 * time.Second)
}
