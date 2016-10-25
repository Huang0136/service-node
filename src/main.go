package main

import (
	"logs"
	_ "service"
	"time"

	_ "github.com/coreos/etcd/clientv3"
	_ "golang.org/x/net/context"
)

var Count int = 0

func main() {

	for {

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
