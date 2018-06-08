package redisutil

import (
	"testing"
)

const redisServerURL = "redis://:123456@192.168.8.175:7001/0"

func TestRedisClient_Ping(t *testing.T) {
	redisClient := GetRedisClient(redisServerURL)
	val, err := redisClient.Ping()
	if err != nil{
		t.Error(err)
	}else{
		t.Log(val)
	}
}