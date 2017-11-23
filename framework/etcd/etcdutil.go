package etcdutil

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

func NewEtcdClientV3(endPoints []string, dialTimeout time.Duration) (*clientv3.Client, error){
	if endPoints == nil {
		endPoints = make([]string,0)
		endPoints = append(endPoints, "127.0.0.1:2379")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func BasicEtcdOperate(){
	
}
