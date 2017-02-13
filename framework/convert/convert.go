package convert

import (
	"errors"
	"math/big"
	"strconv"
	"time"
)

//convert string to []byte
func String2Bytes(val string) []byte {
	return []byte(val)
}

//convert string to int
func String2Int(val string) (int, error) {
	return strconv.Atoi(val)
}

//convert int to string
func Int2String(val int) string {
	return strconv.Itoa(val)
}

func String2Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

func Int642String(val int64) string {
	return strconv.FormatInt(val, 10)
}

func NSToTime(ns int64) (time.Time, error) {
	if ns <= 0 {
		return time.Time{}, errors.New("ns is err")
	}
	bigNS := big.NewInt(ns)
	return time.Unix(ns/1e9, int64(bigNS.Mod(bigNS, big.NewInt(1e9)).Uint64())), nil
}
