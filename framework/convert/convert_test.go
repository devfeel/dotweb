package convert

import (
	"testing"
	"time"
	"github.com/devfeel/dotweb/test"
)

//功能测试

func Test_String2Bytes_1(t *testing.T) {
	str := "0123456789"
	b := String2Bytes(str)
	t.Log(str, " String to Byte: ", b)
	excepted:=[]byte{48,49,50,51,52,53,54,55,56,57}
	test.Equal(t,excepted,b)
}

func Test_String2Int_1(t *testing.T) {
	str := "1234567890"
	b, e := String2Int(str)

	t.Log(str, " String to Int: ", b)
	test.Nil(t,e)
	test.Equal(t,1234567890,b)
}

func Test_String2Int_2(t *testing.T) {
	str := "1234567890ssss"
	b, e := String2Int(str)

	t.Log(str, " String to Int: ", b)
	test.NotNil(t,e)
	test.Equal(t,0,b)
}

func Test_Int2String_1(t *testing.T) {
	vint := 9876543210
	s := Int2String(vint)
	t.Log(vint, "Int to String: ", s)
	test.Equal(t,"9876543210",s)
}

//String2Int64
func Test_String2Int64_1(t *testing.T) {
	str := "0200000010"
	b, e := String2Int64(str)

	t.Log(str, "String to Int64: ", b)
	test.Nil(t,e)
	test.Equal(t,int64(200000010),b)
}

//String2Int64
func Test_String2Int64_2(t *testing.T) {
	str := "a0200000010"
	b, e := String2Int64(str)

	t.Log(str, "String to Int64: ", b)
	test.NotNil(t,e)
	test.Equal(t,int64(0),b)
}

//Int642String
func Test_Int642String_1(t *testing.T) {
	var vint int64 = 1 << 62
	s := Int642String(vint)
	t.Log(vint, "Int64 to String: ", s)
	test.Equal(t,"4611686018427387904",s)
}

func Test_Int642String_2(t *testing.T) {
	var vint int64 = 1 << 62 >> 4
	s := Int642String(vint)
	t.Log(vint, "Int64 to String: ", s)

	test.Equal(t,"288230376151711744",s)
}

//NSToTime
func Test_NSToTime_1(t *testing.T) {
	now := time.Now().UnixNano()
	b, e := NSToTime(now)
	test.Nil(t,e)
	t.Log(now, "NSToTime: ", b)
}

//NSToTime
func Test_NSToTime_2(t *testing.T) {
	now := time.Now().Unix()
	b, e := NSToTime(now)
	test.Nil(t,e)
	t.Log(now, "NSToTime: ", b)
}
