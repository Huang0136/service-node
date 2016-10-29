package main

import (
	_ "constants"
	"logs"
	"service"
	"time"
)

func main() {

	for {
		logs.MyDebugLog.Println("server node working...")
		time.Sleep(time.Minute)
	}
}

func init() {
	go service.RegisterRpc()

}
