package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	ctx1 "golang.org/x/net/context"
)

var DebugLogger *log.Logger

func main() {
	c := clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   []string{"localhost:2379"},
	}



	cli, err := clientv3.New(c)
	defer cli.Close()
	CheckError(err)

	fmt.Println("创建etcd ClientV3成功,", cli)

	resp, err := cli.Put(ctx1.TODO(), "/servers1", "所有节点")
	CheckError(err)

	if resp.PrevKv != nil {
		fmt.Println("key:", string(resp.PrevKv.Key), ",value:", string(resp.PrevKv.Value))
	} else {
		fmt.Println("prevKV is null")
	}

	ctxTemp, cancel := ctx1.WithTimeout(ctx1.Background(), 3)
	resp1, err := cli.Get(ctxTemp, "/servers1")
	cancel()
	CheckError(err)

	for _, ev := range resp1.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	DebugLogger.Println("server node startup")
}

func init() {
	fileLog, err := os.Create("node.log")
	if err != nil {
		fmt.Println("create log file error:", err)
	}

	DebugLogger = log.New(fileLog, "[DEBUG] ", log.LstdFlags)
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}
