package artisan

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5bin(bin []byte) string {
	h := md5.New()
	h.Write(bin)
	return hex.EncodeToString(h.Sum(nil))
}
