package utils

import (
	"bytes"
	"constants"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

// 根据 ip:port:timestamp 生成Token
func CreateToken() string {
	var b bytes.Buffer
	b.WriteString(constants.Configs["serverNode.ip"]) // IP
	b.WriteString(":")
	b.WriteString(constants.Configs["serverNode.port"]) // PORT
	b.WriteString(":")
	b.WriteString(strconv.Itoa(int(time.Now().UnixNano()))) // timestamp

	md5Hash := md5.New()
	io.WriteString(md5Hash, b.String())

	return fmt.Sprintf("%x", md5Hash.Sum(nil)) // 转换成16进制
}
