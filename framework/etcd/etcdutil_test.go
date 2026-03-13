package etcdutil

import (
	"testing"
	"fmt"
	"os"
	"time"
)

var cli *EtcdClient

var err error

func init() {
	cli, err = NewEtcdClient(10*time.Second)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
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
