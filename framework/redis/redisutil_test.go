package redisutil

import (
	"testing"
)

// redisAvailable indicates if Redis server is available for testing
var redisAvailable bool

func init() {
	// Try to connect to Redis at init time
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	_, err := client.Ping()
	redisAvailable = (err == nil)
}

// skipIfNoRedis skips the test if Redis is not available
func skipIfNoRedis(t *testing.T) {
	if !redisAvailable {
		t.Skip("Redis server not available, skipping test")
	}
}

// TestRedisClient_GetDefaultRedisClient tests GetDefaultRedisClient
func TestRedisClient_GetDefaultRedisClient(t *testing.T) {
	// This test doesn't need Redis connection, it just creates a client
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Error("GetDefaultRedisClient returned nil")
	}
}

// TestRedisClient_GetRedisClient tests GetRedisClient with custom pool settings
func TestRedisClient_GetRedisClient(t *testing.T) {
	// This test doesn't need Redis connection
	client := GetRedisClient("redis://localhost:6379/0", 5, 10)
	if client == nil {
		t.Error("GetRedisClient returned nil")
	}

	// Test with zero values (should use defaults)
	client2 := GetRedisClient("redis://localhost:6379/0", 0, 0)
	if client2 == nil {
		t.Error("GetRedisClient with zero values returned nil")
	}
}

// TestRedisClient_Get tests Get operation
func TestRedisClient_Get(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	_, err := client.Get("nonexistent_key_test")
	if err != nil && err.Error() != "redis: nil" {
		t.Logf("Get non-existent key error (expected): %v", err)
	}
}

// TestRedisClient_Set tests Set operation
func TestRedisClient_Set(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_set_key"
	val := "test_value"
	_, err := client.Set(key, val)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	client.Del(key)
}

// TestRedisClient_SetAndGet tests Set followed by Get
func TestRedisClient_SetAndGet(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_setget_key"
	val := "test_value_123"
	_, err := client.Set(key, val)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	got, err := client.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != val {
		t.Errorf("Get returned wrong value: got %s, want %s", got, val)
	}
	client.Del(key)
}

// TestRedisClient_Del tests Del operation
func TestRedisClient_Del(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_del_key"
	client.Set(key, "value")
	_, err := client.Del(key)
	if err != nil {
		t.Errorf("Del failed: %v", err)
	}
	_, err = client.Get(key)
	if err == nil {
		t.Error("Key still exists after Del")
	}
}

// TestRedisClient_Exists tests Exists operation
func TestRedisClient_Exists(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_exists_key"
	exists, err := client.Exists(key)
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Non-existent key should not exist")
	}
	client.Set(key, "value")
	exists, err = client.Exists(key)
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Key should exist after Set")
	}
	client.Del(key)
}

// TestRedisClient_INCR tests INCR operation
func TestRedisClient_INCR(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_incr_key"
	client.Del(key)
	val, err := client.INCR(key)
	if err != nil {
		t.Errorf("INCR failed: %v", err)
	}
	if val != 1 {
		t.Errorf("INCR returned wrong value: got %d, want 1", val)
	}
	val, err = client.INCR(key)
	if err != nil {
		t.Errorf("INCR failed: %v", err)
	}
	if val != 2 {
		t.Errorf("INCR returned wrong value: got %d, want 2", val)
	}
	client.Del(key)
}

// TestRedisClient_DECR tests DECR operation
func TestRedisClient_DECR(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_decr_key"
	client.Del(key)
	val, err := client.DECR(key)
	if err != nil {
		t.Errorf("DECR failed: %v", err)
	}
	if val != -1 {
		t.Errorf("DECR returned wrong value: got %d, want -1", val)
	}
	client.Del(key)
}

// TestRedisClient_Expire tests Expire operation
func TestRedisClient_Expire(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_expire_key"
	client.Set(key, "value")
	_, err := client.Expire(key, 10)
	if err != nil {
		t.Errorf("Expire failed: %v", err)
	}
	client.Del(key)
}

// TestRedisClient_SetWithExpire tests SetWithExpire operation
func TestRedisClient_SetWithExpire(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_setexpire_key"
	_, err := client.SetWithExpire(key, "value", 10)
	if err != nil {
		t.Errorf("SetWithExpire failed: %v", err)
	}
	got, err := client.Get(key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if got != "value" {
		t.Errorf("Get returned wrong value: got %s, want value", got)
	}
	client.Del(key)
}

// TestRedisClient_Ping tests Ping operation
func TestRedisClient_Ping(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	pong, err := client.Ping()
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}
	if pong != "PONG" {
		t.Errorf("Ping returned wrong response: got %s, want PONG", pong)
	}
}

// TestRedisClient_GetConn tests GetConn operation
func TestRedisClient_GetConn(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	conn := client.GetConn()
	if conn == nil {
		t.Error("GetConn returned nil")
		return
	}
	// conn is a connWrapper for backwards compatibility
	t.Log("GetConn returned connection wrapper")
}

// TestRedisClient_MultipleClients tests multiple client instances
func TestRedisClient_MultipleClients(t *testing.T) {
	// This test doesn't need Redis connection
	url := "redis://localhost:6379/0"
	client1 := GetDefaultRedisClient(url)
	client2 := GetDefaultRedisClient(url)
	if client1 != client2 {
		t.Error("GetDefaultRedisClient should return cached instance")
	}
	client3 := GetRedisClient(url, 5, 10)
	client4 := GetRedisClient(url, 5, 10)
	if client3 != client4 {
		t.Error("GetRedisClient should return cached instance for same settings")
	}
}

// TestRedisClient_HashOperations tests HSet, HGet, HGetAll, HDel
func TestRedisClient_HashOperations(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_hash_key"
	client.Del(key)
	err := client.HSet(key, "field1", "value1")
	if err != nil {
		t.Errorf("HSet failed: %v", err)
	}
	val, err := client.HGet(key, "field1")
	if err != nil {
		t.Errorf("HGet failed: %v", err)
	}
	if val != "value1" {
		t.Errorf("HGet returned wrong value: got %s, want value1", val)
	}
	all, err := client.HGetAll(key)
	if err != nil {
		t.Errorf("HGetAll failed: %v", err)
	}
	if all["field1"] != "value1" {
		t.Errorf("HGetAll returned wrong value: got %s, want value1", all["field1"])
	}
	_, err = client.HDel(key, "field1")
	if err != nil {
		t.Errorf("HDel failed: %v", err)
	}
	client.Del(key)
}

// TestRedisClient_ListOperations tests LPush, RPush, LRange, LPop, RPop
func TestRedisClient_ListOperations(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_list_key"
	client.Del(key)
	count, err := client.LPush(key, "value1")
	if err != nil {
		t.Errorf("LPush failed: %v", err)
	}
	if count != 1 {
		t.Errorf("LPush returned wrong count: got %d, want 1", count)
	}
	count, err = client.RPush(key, "value2")
	if err != nil {
		t.Errorf("RPush failed: %v", err)
	}
	if count != 2 {
		t.Errorf("RPush returned wrong count: got %d, want 2", count)
	}
	vals, err := client.LRange(key, 0, -1)
	if err != nil {
		t.Errorf("LRange failed: %v", err)
	}
	if len(vals) != 2 {
		t.Errorf("LRange returned wrong count: got %d, want 2", len(vals))
	}
	val, err := client.LPop(key)
	if err != nil {
		t.Errorf("LPop failed: %v", err)
	}
	t.Logf("LPop: %s", val)
	val, err = client.RPop(key)
	if err != nil {
		t.Errorf("RPop failed: %v", err)
	}
	t.Logf("RPop: %s", val)
	client.Del(key)
}

// TestRedisClient_SetOperations tests SAdd, SMembers, SIsMember, SRem
func TestRedisClient_SetOperations(t *testing.T) {
	skipIfNoRedis(t)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	key := "test_set_key"
	client.Del(key)
	count, err := client.SAdd(key, "member1", "member2")
	if err != nil {
		t.Errorf("SAdd failed: %v", err)
	}
	if count != 2 {
		t.Errorf("SAdd returned wrong count: got %d, want 2", count)
	}
	members, err := client.SMembers(key)
	if err != nil {
		t.Errorf("SMembers failed: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("SMembers returned wrong count: got %d, want 2", len(members))
	}
	isMember, err := client.SIsMember(key, "member1")
	if err != nil {
		t.Errorf("SIsMember failed: %v", err)
	}
	if !isMember {
		t.Error("SIsMember returned false for existing member")
	}
	count, err = client.SRem(key, "member1")
	if err != nil {
		t.Errorf("SRem failed: %v", err)
	}
	if count != 1 {
		t.Errorf("SRem returned wrong count: got %d, want 1", count)
	}
	client.Del(key)
}
