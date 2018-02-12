package convert

import (
	"errors"
	"math/big"
	"strconv"
	"time"
)

// String2Bytes convert string to []byte
func String2Bytes(val string) []byte {
	return []byte(val)
}

// String2Int convert string to int
func String2Int(val string) (int, error) {
	return strconv.Atoi(val)
}

// Int2String convert int to string
func Int2String(val int) string {
	return strconv.Itoa(val)
}

// String2Int64 convert string to int64
func String2Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

// Int642String convert int64 to string
func Int642String(val int64) string {
	return strconv.FormatInt(val, 10)
}

// String2UInt64 convert string to uint64
func String2UInt64(val string) (uint64, error) {
	return strconv.ParseUint(val, 10, 64)
}

// UInt642String convert uint64 to string
func UInt642String(val uint64) string{
	return strconv.FormatUint(val, 10)
}

// NSToTime convert ns to time.Time
func NSToTime(ns int64) (time.Time, error) {
	if ns <= 0 {
		return time.Time{}, errors.New("ns is err")
	}
	bigNS := big.NewInt(ns)
	return time.Unix(ns/1e9, int64(bigNS.Mod(bigNS, big.NewInt(1e9)).Uint64())), nil
}
