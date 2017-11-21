package uuid

import (
	"github.com/devfeel/dotweb/test"
	"testing"
)

// Test_GetUUID_V1_32 test uuid with v1 and return 32 len string
func Test_GetUUID_V1_32(t *testing.T) {
	uuid := NewV1().String32()
	t.Log("GetUUID:", uuid)
	test.Equal(t, 32, len(uuid))
}

// Test_GetUUID_V1 test uuid with v1 and return 36 len string
func Test_GetUUID_V1(t *testing.T) {
	uuid := NewV1().String()
	t.Log("GetUUID:", uuid)
	test.Equal(t, 36, len(uuid))
}

func Benchmark_GetUUID_V1_32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV1().String32()
	}
}

// Test_GetUUID_V4_32 test uuid with v1 and return 32 len string
func Test_GetUUID_V4_32(t *testing.T) {
	uuid := NewV4().String32()
	t.Log("GetUUID:", uuid)
	test.Equal(t, 32, len(uuid))
}

// Test_GetUUID_V4 test uuid with v1 and return 36 len string
func Test_GetUUID_V4(t *testing.T) {
	uuid := NewV4().String()
	t.Log("GetUUID:", uuid)
	test.Equal(t, 36, len(uuid))
}
func Benchmark_GetUUID_V4_32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV4().String32()
	}
}
