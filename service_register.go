package go_service_mgr

import (
	v3 "go.etcd.io/etcd/client/v3"
	"time"
)

type RegSvcClient struct {
	client *v3.Client
}

func CreateRegSvcClient(endpoints []string, dialTimeout uint64) (*RegSvcClient, error) {
	config := v3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
	}
	client, err := v3.New(config)
	if err != nil {
		return nil, err
	}
	return &RegSvcClient{
		client: client,
	}, nil
}
