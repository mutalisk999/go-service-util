package go_service_util

import (
	"context"
	"crypto/tls"
	v3 "go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

type RegSvcClient struct {
	client *v3.Client
}

type RegSvcAuth struct {
	Username string
	Password string
}

func CreateRegSvcClient(endpoints []string, tlsConfig *tls.Config, authConfig *RegSvcAuth, dialTimeout uint64) (*RegSvcClient, error) {
	config := v3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
		TLS:         tlsConfig,
	}
	if authConfig != nil {
		config.Username = authConfig.Username
		config.Password = authConfig.Password
	}

	client, err := v3.New(config)
	if err != nil {
		return nil, err
	}
	return &RegSvcClient{
		client: client,
	}, nil
}

func (c *RegSvcClient) DisposeRegSvcClient() error {
	err := c.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *RegSvcClient) RegisterService(ctxOpTimeout uint64, serviceTTL uint64,
	keyPrefix, id, name, address string, logger *log.Logger) (context.CancelFunc, error) {

	if keyPrefix == "" {
		// default key prefix
		keyPrefix = "/etcd_services"
	}

	if serviceTTL == 0 {
		// default lease ttl
		serviceTTL = 30
	}

	if ctxOpTimeout == 0 {
		// default operation timeout
		ctxOpTimeout = 5
	}

	kv := v3.NewKV(c.client)
	lease := v3.NewLease(c.client)

	ctx, _ := ServiceContextWithTimeout(ctxOpTimeout)
	grantResp, err := lease.Grant(ctx, int64(serviceTTL))
	if err != nil {
		return nil, err
	}

	ctx, _ = ServiceContextWithTimeout(ctxOpTimeout)
	key := strings.Join([]string{keyPrefix, id, name}, "/")
	_, err = kv.Put(ctx, key, address, v3.WithLease(grantResp.ID))
	if err != nil {
		return nil, err
	}

	ctx, cbCancel := ServiceContextWithCancel()
	keepAliveRespChan, err := lease.KeepAlive(ctx, grantResp.ID)
	if err != nil {
		return nil, err
	}

	go keepAliveCallBack(keepAliveRespChan, ctx, key, logger)
	return cbCancel, nil
}

func keepAliveCallBack(keepAliveRespChan <-chan *v3.LeaseKeepAliveResponse, ctx context.Context, key string, logger *log.Logger) {
	for {
		select {
		case ret := <-keepAliveRespChan:
			if ret != nil && logger != nil {
				logger.Printf("Keep Alive Service [%s] at %s", key, time.Now().String())
			}
		case <-ctx.Done():
			if logger != nil {
				logger.Printf("[Cancel] Keep Alive Service [%s] at %s", key, time.Now().String())
			}
			return
		}
	}
}
