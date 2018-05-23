package redisutil

import (
	"testing"
)

const redisServerUrl = "redis://192.168.8.175:6379/0"

func TestRedisClient_Ping(t *testing.T) {
	redisClient := GetRedisClient(redisServerUrl)
	val, err := redisClient.Ping()
	if err != nil{
		t.Error(err)
	}else{
		t.Log(val)
	}
}