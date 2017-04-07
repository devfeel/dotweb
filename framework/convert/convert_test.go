package convert

import (
	"testing"
	"time"
)

//功能测试

func Test_String2Bytes_1(t *testing.T) {
	str := "0123456789"
	b := String2Bytes(str)
	t.Log(str, " String to Byte: ", b)
}

func Test_String2Int_1(t *testing.T) {
	str := "1234567890"
	b, e := String2Int(str)
	t.Error(e)
	t.Log(str, " String to Int: ", b)
}

func Test_String2Int_2(t *testing.T) {
	str := "1234567890ssss"
	b, e := String2Int(str)
	t.Error(e)
	t.Log(str, " String to Int: ", b)
}

func Test_Int2String_1(t *testing.T) {
	vint := 9876543210
	s := Int2String(vint)
	t.Log(vint, "Int to String: ", s)
}

//String2Int64
func Test_String2Int64_1(t *testing.T) {
	str := "0200000010"
	b, e := String2Int64(str)
	t.Error(e)
	t.Log(str, "String to Int64: ", b)
}

//String2Int64
func Test_String2Int64_2(t *testing.T) {
	str := "a0200000010"
	b, e := String2Int64(str)
	t.Error(e)
	t.Log(str, "String to Int64: ", b)
}

//Int642String
func Test_Int642String_1(t *testing.T) {
	var vint int64 = 1 << 62
	s := Int642String(vint)
	t.Log(vint, "Int64 to String: ", s)
}

func Test_Int642String_2(t *testing.T) {
	var vint int64 = 1 << 62 >> 4
	s := Int642String(vint)
	t.Log(vint, "Int64 to String: ", s)
}

//NSToTime
func Test_NSToTime_1(t *testing.T) {
	now := time.Now().UnixNano()
	b, e := NSToTime(now)
	t.Error(e)
	t.Log(now, "NSToTime: ", b)
}

//NSToTime
func Test_NSToTime_2(t *testing.T) {
	now := time.Now().Unix()
	b, e := NSToTime(now)
	t.Error(e)
	t.Log(now, "NSToTime: ", b)
}
