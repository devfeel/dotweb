package etcdutil

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"context"
)
//etcd util only support go1.9+ !!!
type EtcdClient struct {
	client *clientv3.Client
}

func NewEtcdClient(dialTimeout time.Duration, endPoints ...string) (*EtcdClient, error) {
	etcdClient := new(EtcdClient)

	if endPoints == nil {
		endPoints = make([]string, 0)
		endPoints = append(endPoints, "127.0.0.1:2379")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		return nil, err
	}
	etcdClient.client = cli
	return etcdClient, nil
}

func (e *EtcdClient) SimpleGet(key string) (string, error) {
	getResp, err := e.client.Get(context.TODO(), key)
	if err != nil {
		return "", err
	}
	return string(getResp.Kvs[0].Value), nil
}

func (e *EtcdClient) SimplePut(key, value string) error {
	_, err := e.client.Put(context.TODO(), key, value)
	return err
}

func (e *EtcdClient) SimpleDelete(key string) (int64, error) {
	delResp, err := e.client.Delete(context.TODO(), key)
	return delResp.Deleted, err
}

func (e *EtcdClient) SimpleWatch(key string) clientv3.WatchChan {
	return e.client.Watch(context.Background(), key)
}

// etcd v3 not has directory construct,
// by watch a prefix of key to implement watch a directorys in v2
func (e *EtcdClient) PrefixWatch(prefix string) clientv3.WatchChan {
	return e.client.Watch(context.Background(), prefix)
}

func (e *EtcdClient) Close() error {
	err := e.client.Close()
	if err != nil {
		return err
	}
	return nil
}
