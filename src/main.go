package main

import (
	_ "constants"
	"logs"
	_ "register"
	_ "service"
	"time"
)

func main() {

	for {
		logs.MyDebugLog.Println("server node working...")
		time.Sleep(time.Minute)
	}
}

func init() {
	// 初始化
	//	constants.init()

	// 注册节点和服务接口到etcd
	//	register.RegisterToEtcd()

	// 开启rpc监听
	//	go service.RegisterRpc()

}
