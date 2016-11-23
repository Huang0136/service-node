package beans

import (
	"time"
)

// 用户
type User struct {
	Id     int
	Name   string // 用户名
	Status int    // 状态
	Time   time.Time
}
