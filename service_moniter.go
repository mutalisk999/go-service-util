package go_service_util

import (
	"context"
	"crypto/tls"
	"go.etcd.io/etcd/api/v3/mvccpb"
	v3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/syncmap"
	"log"
	"time"
)

type MonSvcClient struct {
	client *v3.Client
	svcMap syncmap.Map
}

type MonSvcAuth struct {
	Username string
	Password string
}

func CreateMonSvcClient(endpoints []string, tlsConfig *tls.Config, authConfig *MonSvcAuth, dialTimeout uint64) (*MonSvcClient, error) {
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
	return &MonSvcClient{
		client: client,
	}, nil
}

func (c *MonSvcClient) DisposeMonSvcClient() error {
	err := c.client.Close()
	if err != nil {
		return err
	}
	c.svcMap = syncmap.Map{}
	return nil
}

func (c *MonSvcClient) GetService() map[string]string {
	svcMapRet := make(map[string]string)
	c.svcMap.Range(func(k interface{}, v interface{}) bool {
		svcMapRet[k.(string)] = v.(string)
		return true
	})
	return svcMapRet
}

func (c *MonSvcClient) MonitorService(ctxOpTimeout uint64, keyPrefix string, logger *log.Logger,
	putCallBack func(string, string), deleteCallBack func(string)) (context.CancelFunc, error) {
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

	for _, kv := range getResp.Kvs {
		v := kv.Value
		if v != nil {
			c.svcMap.Store(string(kv.Key), string(kv.Value))
		}
	}

	ctx, cbCancel := ServiceContextWithCancel()
	watchChan := c.client.Watch(ctx, keyPrefix, v3.WithPrefix())

	go c.watcherCallBack(watchChan, ctx, keyPrefix, logger, putCallBack, deleteCallBack)
	return cbCancel, nil
}

func (c *MonSvcClient) watcherCallBack(watchChan v3.WatchChan, ctx context.Context, keyPrefix string, logger *log.Logger,
	putCallBack func(string, string), deleteCallBack func(string)) {
	for {
		select {
		case ret := <-watchChan:
			if ret.Events == nil {
				continue
			}
			for _, ev := range ret.Events {
				switch ev.Type {
				case mvccpb.PUT:
					c.svcMap.Store(string(ev.Kv.Key), string(ev.Kv.Value))
					if putCallBack != nil {
						putCallBack(string(ev.Kv.Key), string(ev.Kv.Value))
					}
					if logger != nil {
						logger.Printf("Service Watcher Put [%s] | [%s] at %s",
							string(ev.Kv.Key), string(ev.Kv.Value), time.Now().String())
					}
				case mvccpb.DELETE:
					c.svcMap.Delete(string(ev.Kv.Key))
					if deleteCallBack != nil {
						deleteCallBack(string(ev.Kv.Key))
					}
					if logger != nil {
						logger.Printf("Service Watcher Delete [%s] at %s",
							string(ev.Kv.Key), time.Now().String())
					}
				}
			}
		case <-ctx.Done():
			if logger != nil {
				logger.Printf("Cancel Service Watcher [%s] at %s", keyPrefix, time.Now().String())
			}
			return
		}
	}
}
