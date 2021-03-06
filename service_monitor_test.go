package go_service_util

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestCreateMonSvcClient(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	svcCli, err := CreateMonSvcClient([]string{"http://127.0.0.1:12379", "http://127.0.0.1:22379", "http://127.0.0.1:32379"}, nil, nil, 3)
	if err != nil {
		logger.Println("CreateRegSvcClient, err:", err)
		return
	}
	defer svcCli.DisposeMonSvcClient()

	cbCancel, err := svcCli.MonitorService(3, "", logger, nil, nil)
	if err != nil {
		logger.Println("MonitorService, err:", err)
		return
	}

	time.Sleep(100 * time.Second)
	logger.Println("svcCli.GetService: ", svcCli.GetService())
	time.Sleep(300 * time.Second)
	logger.Println("svcCli.GetService: ", svcCli.GetService())
	cbCancel()
	time.Sleep(10 * time.Second)
}
