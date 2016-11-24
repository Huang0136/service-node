package beans

import (
	"time"
)

// 用户实体
type User struct {
	Id     int       `json:"id"`
	Name   string    `json:"name"`   // 用户名
	Status int       `json:"status"` // 状态
	Time   time.Time `json:"register_date"`
}
