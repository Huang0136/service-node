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
	RegisterToEtcd1()
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
		DialTimeout: 1 * time.Minute,
	})

	if err != nil {
		log.Fatal("创建etcd clientv3失败:", err)
	}
	defer cli.Close()

	/*
		delResp, err := cli.Delete(context.TODO(), "", clientv3.WithPrefix())
		if err != nil {
			log.Fatalln("delete error:", err)
		}
		fmt.Println("删除成功", delResp)
	*/

	// 10秒后设置值
	time.Sleep(5 * time.Second)

	// 设置超时
	gresp, err := cli.Grant(context.Background(), 20)
	if err != nil {
		log.Fatalln("grant失败:", err, ",", gresp)
	}

	// put
	putResp, err := cli.Put(context.Background(), "servers/127.0.0.1:9090", "这是服务节点server-node,Time:"+time.Now().Format("2006-01-02 15:04:05.999")) //, clientv3.WithLease(gresp.ID)
	if err != nil {
		log.Fatalln("设值失败:", err)
	}
	fmt.Println("服务节点注册成功:", putResp)

	/*
		time.Sleep(3 * time.Second)
		// put
		putResp1, err := cli.Put(context.TODO(), "servers/127.0.0.1:9090/service_001", "服务接口001")
		if err != nil {
			log.Fatalln("设值失败2:", err)
		}
		fmt.Println("服务节点注册成功2:", putResp1)
	*/

	// 注册服务
	// servers/ip:port/URL:MethodType

	logs.MyInfoLog.Println("注册成功!")
}

func RegisterToEtcd1() {
	logs.MyInfoLog.Println("监听Etcd...")

	// 注册节点
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.Configs["etcd.url"]},
		DialTimeout: 1 * time.Hour,
	})

	if err != nil {
		log.Fatal("创建etcd clientv3失败:", err)
	}
	defer cli.Close()

	// watch
	go func() {
		count := 1
	WatchBlock:
		{
			//			fmt.Println("开始监听服务节点", count, time.Now().Format("2006-01-02 15:04:05.999"))
			wc1 := cli.Watch(context.Background(), "servers/127.0.0.1:9090")
			//			fmt.Println("监听到服务节点变化", count, time.Now().Format("2006-01-02 15:04:05.999"))

			/*
				for wcTemp := range wc1 {
					for _, wcEvent := range wcTemp.Events {
						fmt.Printf("1 Type:%s,key:%s,value:%s \n", wcEvent.Type, string(wcEvent.Kv.Key), string(wcEvent.Kv.Value))
					}
				}
			*/

			//			fmt.Println("我想执行1 ", count, " ", time.Now())
			//			watchResp := <-wc1

			//			fmt.Println(watchResp, "==========", watchResp.Header.Revision)
			//			for _, e := range watchResp.Events {
			//				fmt.Printf("Type:%s,key:%s,value:%s ,Time:%s\n", e.Type, e.Kv.Key, e.Kv.Value, time.Now().Format("2006-01-02 15:04:05.999"))
			//			}

			for wresp := range wc1 {
				for _, ev := range wresp.Events {
					fmt.Printf("监听到信息,%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				}
			}
			//			fmt.Println("我想执行2 ", count, " ", time.Now())

			count++
			if count > 100000000 {
				goto bb
			}
			goto WatchBlock
		}

	bb:
		{
			fmt.Println("已然结束!")
		}

	}()

	logs.MyInfoLog.Println("监听Etcd结束!")
}
