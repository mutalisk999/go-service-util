package go_service_util

import (
	v3 "go.etcd.io/etcd/client/v3"
	"time"
)

type MonSvcClient struct {
	client *v3.Client
}

func CreateMonSvcClient(endpoints []string, dialTimeout uint64) (*MonSvcClient, error) {
	config := v3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
	}
	client, err := v3.New(config)
	if err != nil {
		return nil, err
	}
	return &MonSvcClient{
		client: client,
	}, nil
}

func (c *MonSvcClient) DisposeMonSvcClient() error {
	err := c.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *MonSvcClient) GetService(ctxOpTimeout uint64, keyPrefix string) (map[string]string, error) {
	if keyPrefix == "" {
		// default key prefix
		keyPrefix = "/etcd_services"
	}

	if ctxOpTimeout == 0 {
		// default operation timeout
		ctxOpTimeout = 5
	}

	ctx, _ := ServiceContextWithTimeout(ctxOpTimeout)
	getResp, err := c.client.Get(ctx, keyPrefix, v3.WithPrefix())
	if err != nil {
		return nil, nil
	}

	if getResp == nil || getResp.Kvs == nil {
		return nil, nil
	}

	svcMap := make(map[string]string)
	for _, kv := range getResp.Kvs {
		v := kv.Value
		if v != nil {
			svcMap[string(kv.Key)] = string(kv.Value)
		}
	}
	return svcMap, nil
}
