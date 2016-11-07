package register

import (
	"fmt"
	"log"
	"logs"
	"regexp"
	"sync"
	"time"

	"constants"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

var AllHandlers Handlers

func init() {
	// 注册到Etcd
	RegisterToEtcd()

	/*
		AllHandlers = new(Handlers)
		AllHandlers.GetHandlers = make(map[string][]Service)
		AllHandlers.PostHandlers = make(map[string][]Service)
		AllHandlers.PutHandlers = make(map[string][]Service)
		AllHandlers.DeleteHandlers = make(map[string][]Service)
		AllHandlers.GetRegexHandlers = make(map[string][]Service)
		AllHandlers.PostRegexHandlers = make(map[string][]Service)
		AllHandlers.PutRegexHandlers = make(map[string][]Service)
		AllHandlers.DeleteRegexHandlers = make(map[string][]Service)

		serverNodeList = make([]ServerNode)*/

}

// 处理方法
type Handlers struct {
	GetHandlers    map[string][]Service // GET类型精确的处理器
	PostHandlers   map[string][]Service // POST类型精确的处理器
	PutHandlers    map[string][]Service // PUT类型精确的处理器
	DeleteHandlers map[string][]Service // DELETE类型精确的处理器

	GetRegexHandlers    map[string][]Service // GET类型模糊的处理器
	PostRegexHandlers   map[string][]Service // POST类型模糊的处理器
	PutRegexHandlers    map[string][]Service // PUT类型模糊的处理器
	DeleteRegexHandlers map[string][]Service // DELETE类型模糊的处理器

	Lock sync.Mutex // 锁
}

// 服务节点列表
var serverNodeList []ServerNode

// 服务节点
type ServerNode struct {
	IP           string    // IP
	Port         int       // 端口
	Instance     string    // 节点实例名
	Desc         string    // 描述
	Remark       string    // 备注
	RegisterDate time.Time // 注册时间
	Enable       bool      // 是否可访问
}

// 服务接口列表
var serviceList []Service

// 服务接口
type Service struct {
	ServiceName string        // 接口名称
	MethodName  string        // 方法名称
	MethodType  string        // 方法类型
	InParams    string        // 入参
	OutParams   string        // 出参
	Node        ServerNode    // 所属节点
	RegexStr    regexp.Regexp // 正则表达式
}

// 监听注册中心Etcd
func watch() {

}

// 转换服务节点以及服务接口信息
func transform() {
	AllHandlers.Lock.Lock()
	defer AllHandlers.Lock.Unlock()

	// transform

}

// 注册节点信息到Etcd(服务节点、服务接口等信息)
func RegisterToEtcd() {
	logs.MyInfoLog.Println("注册到Etcd...")

	// 注册节点
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.Configs["etcd.url"]},
		DialTimeout: 2 * time.Minute,
	})

	if err != nil {
		log.Fatal("创建etcd clientv3失败:", err)
	}
	defer cli.Close()
	fmt.Println("创建etcd clientv3成功:", cli)

	/*
		go func() {
			// watch
			fmt.Println("开始监听")
			wc1 := cli.Watch(context.TODO(), "servers/127.0.0.1:9090/service001/t1", clientv3.WithPrefix())
			fmt.Println("监听到结果")

			for wcTemp := range wc1 {
				for _, wcEvent := range wcTemp.Events {
					fmt.Printf("1 Type:%s,key:%s,value:%s \n", wcEvent.Type, string(wcEvent.Kv.Key), string(wcEvent.Kv.Value))
				}
			}
		}()

		kv := clientv3.NewKV(cli)

		// put
		ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
		_, err = kv.Put(ctx1, "servers/127.0.0.1:9090/service001/t1", "服务接口001/t1")
		cancel1()
		if err != nil {
			log.Fatalln(err)
		}

		resp1, err := cli.Get(context.TODO(), "servers", clientv3.WithPrefix())
		if err != nil {
			fmt.Println(err)
		}

		for kkkvvv := range resp1.Kvs {
			fmt.Println("", kkkvvv)
		}

		// watch
		fmt.Println("开始监听2")
		watchChan := cli.Watch(context.TODO(), "servers/127.0.0.1:9090/service001/t1", clientv3.WithPrefix())
		fmt.Println("监听到结果2")
		for wcTemp := range watchChan {
			for _, wcEvent := range wcTemp.Events {
				fmt.Printf("2 Type:%s,key:%s,value:%s \n", wcEvent.Type, string(wcEvent.Kv.Key), string(wcEvent.Kv.Value))
			}
		}*/

	ctx, cancel := context.WithTimeout(context.Background(), 105*time.Second)
	resp, err := cli.Get(ctx, "s", clientv3.WithPrefix()) // /127.0.0.1:9090/service001
	cancel()

	if err != nil {
		log.Fatal("get操作失败:", err)
	}

	fmt.Println("返回结果集:", resp)
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}

	// 注册服务
	// servers/ip:port/URL:MethodType

	logs.MyInfoLog.Println("注册成功!")
}
