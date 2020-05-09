package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Hash(context string) string {
	h := md5.New()
	h.Write([]byte(context))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
