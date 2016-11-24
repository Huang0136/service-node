package utils

import (
	"crypto/sha512"
	"fmt"
	"io"
)

// 密码加密方法
// 加密规则:sha512(sha512(password),salt)
func EncryptPassword(pw, slat string) (encryptPw string) {
	// 第一次加密 sha512(password)
	hash1 := sha512.New()
	io.WriteString(hash1, pw)

	pwEncrypt1 := fmt.Sprintf("%x", hash1.Sum(nil))

	// 第二次加密 sha512(pw1,salt)
	hash2 := sha512.New()
	io.WriteString(hash2, pwEncrypt1)
	io.WriteString(hash2, slat)

	return fmt.Sprintf("%x", hash2.Sum(nil))
}
