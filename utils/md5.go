package utils

import (
	"crypto/md5"
    "encoding/hex"
)

func MD5String(value string) string {
	hash := md5.Sum([]byte(value))
	return hex.EncodeToString(hash[:])
}
