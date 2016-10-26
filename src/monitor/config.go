package monitor

import (
	"logs"
	"time"
)

func init() {
	go MonitorConfig()
}

func MonitorConfig() {
	for {
		logs.MyDebugLog.Println("Monitor config file ")
		time.Sleep(time.Second)
	}
}
