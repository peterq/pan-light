package util

import (
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
