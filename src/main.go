package main

import (
	_ "constants"
	"logs"
	_ "register"
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
	//	register.RegisterToEtcd()
	go service.RegisterRpc()

}
