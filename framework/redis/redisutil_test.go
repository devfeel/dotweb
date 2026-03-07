package redisutil

import (
	"testing"
)

// TestRedisClient_GetDefaultRedisClient tests GetDefaultRedisClient
func TestRedisClient_GetDefaultRedisClient(t *testing.T) {
	// This test requires a running Redis server
	// Skip if no Redis server is available
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	// Test with invalid URL (should still return a client, connection happens on use)
	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Error("GetDefaultRedisClient returned nil")
	}
}

// TestRedisClient_GetRedisClient tests GetRedisClient with custom pool settings
func TestRedisClient_GetRedisClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

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
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	// Test get non-existent key
	_, err := client.Get("nonexistent_key_test")
	if err != nil && err.Error() != "redigo: nil returned" {
		t.Logf("Get non-existent key error (expected): %v", err)
	}
}

// TestRedisClient_Set tests Set operation
func TestRedisClient_Set(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_set_key"
	val := "test_value"

	_, err := client.Set(key, val)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_SetAndGet tests Set followed by Get
func TestRedisClient_SetAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_setget_key"
	val := "test_value_123"

	// Set
	_, err := client.Set(key, val)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	got, err := client.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got != val {
		t.Errorf("Get returned wrong value: got %s, want %s", got, val)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_Del tests Del operation
func TestRedisClient_Del(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_del_key"

	// Set a key
	client.Set(key, "value")

	// Delete
	_, err := client.Del(key)
	if err != nil {
		t.Errorf("Del failed: %v", err)
	}

	// Verify deleted
	_, err = client.Get(key)
	if err == nil {
		t.Error("Key still exists after Del")
	}
}

// TestRedisClient_Exists tests Exists operation
func TestRedisClient_Exists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_exists_key"

	// Non-existent key
	exists, err := client.Exists(key)
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Non-existent key should not exist")
	}

	// Set key
	client.Set(key, "value")

	// Existing key
	exists, err = client.Exists(key)
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Key should exist after Set")
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_INCR tests INCR operation
func TestRedisClient_INCR(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_incr_key"

	// Clean up first
	client.Del(key)

	// INCR on non-existent key (should start at 0, return 1)
	val, err := client.INCR(key)
	if err != nil {
		t.Errorf("INCR failed: %v", err)
	}
	if val != 1 {
		t.Errorf("INCR returned wrong value: got %d, want 1", val)
	}

	// INCR again
	val, err = client.INCR(key)
	if err != nil {
		t.Errorf("INCR failed: %v", err)
	}
	if val != 2 {
		t.Errorf("INCR returned wrong value: got %d, want 2", val)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_DECR tests DECR operation
func TestRedisClient_DECR(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_decr_key"

	// Clean up first
	client.Del(key)

	// DECR on non-existent key (should start at 0, return -1)
	val, err := client.DECR(key)
	if err != nil {
		t.Errorf("DECR failed: %v", err)
	}
	if val != -1 {
		t.Errorf("DECR returned wrong value: got %d, want -1", val)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_Expire tests Expire operation
func TestRedisClient_Expire(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_expire_key"

	// Set a key
	client.Set(key, "value")

	// Set expire
	_, err := client.Expire(key, 10)
	if err != nil {
		t.Errorf("Expire failed: %v", err)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_SetWithExpire tests SetWithExpire operation
func TestRedisClient_SetWithExpire(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_setexpire_key"

	// Set with expire
	_, err := client.SetWithExpire(key, "value", 10)
	if err != nil {
		t.Errorf("SetWithExpire failed: %v", err)
	}

	// Verify set
	got, err := client.Get(key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if got != "value" {
		t.Errorf("Get returned wrong value: got %s, want value", got)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_SetNX tests SetNX operation
func TestRedisClient_SetNX(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_setnx_key"

	// Clean up first
	client.Del(key)

	// SetNX on non-existent key (should succeed)
	val, err := client.SetNX(key, "value1")
	if err != nil {
		t.Errorf("SetNX failed: %v", err)
	}
	t.Logf("SetNX result: %v", val)

	// SetNX on existing key (should fail)
	val2, err := client.SetNX(key, "value2")
	if err != nil {
		t.Errorf("SetNX failed: %v", err)
	}
	t.Logf("SetNX result on existing key: %v", val2)

	// Cleanup
	client.Del(key)
}

// TestRedisClient_HashOperations tests HSet, HGet, HGetAll, HDel
func TestRedisClient_HashOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_hash_key"

	// Clean up
	client.Del(key)

	// HSet
	err := client.HSet(key, "field1", "value1")
	if err != nil {
		t.Errorf("HSet failed: %v", err)
	}

	// HGet
	val, err := client.HGet(key, "field1")
	if err != nil {
		t.Errorf("HGet failed: %v", err)
	}
	if val != "value1" {
		t.Errorf("HGet returned wrong value: got %s, want value1", val)
	}

	// HGetAll
	all, err := client.HGetAll(key)
	if err != nil {
		t.Errorf("HGetAll failed: %v", err)
	}
	if all["field1"] != "value1" {
		t.Errorf("HGetAll returned wrong value: got %s, want value1", all["field1"])
	}

	// HDel
	_, err = client.HDel(key, "field1")
	if err != nil {
		t.Errorf("HDel failed: %v", err)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_ListOperations tests LPush, RPush, LRange, LPop, RPop
func TestRedisClient_ListOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_list_key"

	// Clean up
	client.Del(key)

	// LPush
	count, err := client.LPush(key, "value1")
	if err != nil {
		t.Errorf("LPush failed: %v", err)
	}
	if count != 1 {
		t.Errorf("LPush returned wrong count: got %d, want 1", count)
	}

	// RPush
	count, err = client.RPush(key, "value2")
	if err != nil {
		t.Errorf("RPush failed: %v", err)
	}
	if count != 2 {
		t.Errorf("RPush returned wrong count: got %d, want 2", count)
	}

	// LRange
	vals, err := client.LRange(key, 0, -1)
	if err != nil {
		t.Errorf("LRange failed: %v", err)
	}
	if len(vals) != 2 {
		t.Errorf("LRange returned wrong count: got %d, want 2", len(vals))
	}

	// LPop
	val, err := client.LPop(key)
	if err != nil {
		t.Errorf("LPop failed: %v", err)
	}
	t.Logf("LPop: %s", val)

	// RPop
	val, err = client.RPop(key)
	if err != nil {
		t.Errorf("RPop failed: %v", err)
	}
	t.Logf("RPop: %s", val)

	// Cleanup
	client.Del(key)
}

// TestRedisClient_SetOperations tests SAdd, SMembers, SIsMember, SRem
func TestRedisClient_SetOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_set_key"

	// Clean up
	client.Del(key)

	// SAdd
	count, err := client.SAdd(key, "member1", "member2")
	if err != nil {
		t.Errorf("SAdd failed: %v", err)
	}
	if count != 2 {
		t.Errorf("SAdd returned wrong count: got %d, want 2", count)
	}

	// SMembers
	members, err := client.SMembers(key)
	if err != nil {
		t.Errorf("SMembers failed: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("SMembers returned wrong count: got %d, want 2", len(members))
	}

	// SIsMember
	isMember, err := client.SIsMember(key, "member1")
	if err != nil {
		t.Errorf("SIsMember failed: %v", err)
	}
	if !isMember {
		t.Error("SIsMember returned false for existing member")
	}

	// SRem
	count, err = client.SRem(key, "member1")
	if err != nil {
		t.Errorf("SRem failed: %v", err)
	}
	if count != 1 {
		t.Errorf("SRem returned wrong count: got %d, want 1", count)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_Ping tests Ping operation
func TestRedisClient_Ping(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

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
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	conn := client.GetConn()
	if conn == nil {
		t.Error("GetConn returned nil")
		return
	}
	defer conn.Close()

	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		t.Errorf("Connection PING failed: %v", err)
	}
}

// TestRedisClient_Append tests Append operation
func TestRedisClient_Append(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_append_key"

	// Clean up
	client.Del(key)

	// Append to non-existent key (behaves like Set)
	_, err := client.Append(key, "Hello")
	if err != nil {
		t.Errorf("Append failed: %v", err)
	}

	// Append to existing key
	_, err = client.Append(key, " World")
	if err != nil {
		t.Errorf("Append failed: %v", err)
	}

	// Verify
	got, err := client.Get(key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if got != "Hello World" {
		t.Errorf("Get returned wrong value: got %s, want Hello World", got)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_HIncrBy tests HIncrBy operation
func TestRedisClient_HIncrBy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_hincr_key"

	// Clean up
	client.Del(key)

	// HIncrBy on non-existent field
	val, err := client.HIncrBy(key, "counter", 5)
	if err != nil {
		t.Errorf("HIncrBy failed: %v", err)
	}
	if val != 5 {
		t.Errorf("HIncrBy returned wrong value: got %d, want 5", val)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_LLen tests LLen operation
func TestRedisClient_LLen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_llen_key"

	// Clean up
	client.Del(key)

	// LLen on empty list
	len, err := client.LLen(key)
	if err != nil {
		t.Errorf("LLen failed: %v", err)
	}
	if len != 0 {
		t.Errorf("LLen returned wrong value: got %d, want 0", len)
	}

	// Push some values
	client.LPush(key, "v1", "v2", "v3")

	// LLen on non-empty list
	len, err = client.LLen(key)
	if err != nil {
		t.Errorf("LLen failed: %v", err)
	}
	if len != 3 {
		t.Errorf("LLen returned wrong value: got %d, want 3", len)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_SCard tests SCard operation
func TestRedisClient_SCard(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_scard_key"

	// Clean up
	client.Del(key)

	// SCard on empty set
	count, err := client.SCard(key)
	if err != nil {
		t.Errorf("SCard failed: %v", err)
	}
	if count != 0 {
		t.Errorf("SCard returned wrong value: got %d, want 0", count)
	}

	// Add some members
	client.SAdd(key, "m1", "m2", "m3")

	// SCard on non-empty set
	count, err = client.SCard(key)
	if err != nil {
		t.Errorf("SCard failed: %v", err)
	}
	if count != 3 {
		t.Errorf("SCard returned wrong value: got %d, want 3", count)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_GetObj tests GetObj operation
func TestRedisClient_GetObj(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_getobj_key"

	// Clean up
	client.Del(key)

	// Set a value
	client.Set(key, "test_value")

	// GetObj
	val, err := client.GetObj(key)
	if err != nil {
		t.Errorf("GetObj failed: %v", err)
	}
	t.Logf("GetObj returned: %v (type: %T)", val, val)

	// Cleanup
	client.Del(key)
}

// TestRedisClient_HLen tests HLen operation
func TestRedisClient_HLen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	key := "test_hlen_key"

	// Clean up
	client.Del(key)

	// HSet some fields
	client.HSet(key, "f1", "v1")
	client.HSet(key, "f2", "v2")
	client.HSet(key, "f3", "v3")

	// HLen
	len, err := client.HLen(key)
	if err != nil {
		t.Errorf("HLen failed: %v", err)
	}
	if len != 3 {
		t.Errorf("HLen returned wrong value: got %d, want 3", len)
	}

	// Cleanup
	client.Del(key)
}

// TestRedisClient_DBSize tests DBSize operation
func TestRedisClient_DBSize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := GetDefaultRedisClient("redis://localhost:6379/0")
	if client == nil {
		t.Skip("Redis client is nil, skipping")
	}

	size, err := client.DBSize()
	if err != nil {
		t.Errorf("DBSize failed: %v", err)
	}
	t.Logf("DBSize: %d", size)
}

// TestRedisClient_MultipleClients tests multiple client instances
func TestRedisClient_MultipleClients(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	url := "redis://localhost:6379/0"

	// Get same client multiple times (should return cached instance)
	client1 := GetDefaultRedisClient(url)
	client2 := GetDefaultRedisClient(url)

	if client1 != client2 {
		t.Error("GetDefaultRedisClient should return cached instance")
	}

	// Get with different settings
	client3 := GetRedisClient(url, 5, 10)
	client4 := GetRedisClient(url, 5, 10)

	if client3 != client4 {
		t.Error("GetRedisClient should return cached instance for same settings")
	}
}
