package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

const Key = "hello,world"

// CryptPassword 给密码进行加密
func CryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(Key))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
