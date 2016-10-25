package main

import (
	"logs"
	_ "service"
	"time"

	_ "github.com/coreos/etcd/clientv3"
	_ "golang.org/x/net/context"

	"log"
	"net/http"
	"net/rpc"
	"register"
	"service"

	"github.com/coreos/etcd/clientv3"
)

var Count int = 0

func main() {

	for {

	}

	log.Println("serverTest start...")
	log.Println("正在注册rpc服务")

	myServiceTest := new(service.ServiceTest)
	rpc.Register(myServiceTest)

	// 注册信息到注册中心
	register.RegisterInfo()

	// http方式
	rpc.HandleHTTP()
	err := http.ListenAndServe(":9877", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}

}

func init() {
	go func() {
		for {
			Count++
			logs.MyLogger.Println(Count)
			time.Sleep(time.Second)
		}
	}()
}

func CheckError(err error) {
	if err != nil {
		logs.MyLogger.Fatalln(err)
	}
}
