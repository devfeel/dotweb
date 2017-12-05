package etcdutil

import (
	"testing"
	"fmt"
	"context"
)

var cli *EtcdClient

var err error

func init() {
	cli, err = NewEtcdClient(0)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestEtcdPut(t *testing.T) {
	defer cli.Close()
	putErr := cli.SimplePut("foo", "bar")
	if putErr != nil {
		fmt.Println(putErr)
		return
	}
	fmt.Println("success")
}

func TestEtcdGet(t *testing.T) {
	val, getErr := cli.SimpleGet("foo")
	if getErr != nil {
		fmt.Println(getErr)
	}
	fmt.Println(val)
}
