package etcdutil

import (
	"github.com/coreos/etcd/clientv3"
	"fmt"
	"testing"
	"context"
	"time"
)

var cli *clientv3.Client
var err error

func init() {
	cli, err = NewEtcdClientV3(nil, 10)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestBasicEtcdOperate(t *testing.T) {
	resp, _:= cli.Grant(context.TODO(), 10)
	ctx, _:= context.WithTimeout(context.Background(), 5*time.Second)
	rsp, err:=cli.Put(ctx,"root/game/node-2",`{"addr":"192.168.1.1:9999"}`,clientv3.WithLease(resp.ID))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(rsp.OpResponse().Put().Header.String())
}