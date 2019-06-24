package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func UniqId() string {
	t := time.Now().Unix()
	return strconv.Itoa(int(t)) + strconv.Itoa(int(rand.Float64()*5e3))
}

func First(args ...interface{}) interface{} {
	return args[0]
}

func Second(args ...interface{}) interface{} {
	return args[1]
}

func Md5bin(bin []byte) string {
	h := md5.New()
	h.Write(bin)
	return hex.EncodeToString(h.Sum(nil))
}
