package encryptHelper

import (
	"crypto/md5"
	"fmt"
)

// MD5加密（32位）
func EncryptMD5To32Bit(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
