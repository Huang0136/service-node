package register

import (
	"fmt"
	"logs"
	"regexp"
	"sync"
	"time"

	"constants"

	"net/http"
	"net/rpc"

	"service"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

var AllHandlers Handlers

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

var registerChan chan (bool) = make(chan bool)

func init() {
	// 注册到Etcd
	go registerToEtcd()
	<-registerChan // 阻塞等待注册完成

	registerRpc()

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

// 转换服务节点以及服务接口信息
func transform() {
	AllHandlers.Lock.Lock()
	defer AllHandlers.Lock.Unlock()

	// transform

}

// 注册节点信息到Etcd(服务节点、服务接口等信息)
func registerToEtcd() {
	logs.MyInfoLog.Println("开始注册节点信息到Etcd...")

	// 获取连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.Configs["etcd.url"]},
		DialTimeout: 1 * time.Minute,
	})
	logs.MyErrorLog.CheckFatallnError("创建etcd clientv3失败:", err)
	defer cli.Close()

	// 设置约租
	leaseGrantResp, err := cli.Grant(context.Background(), 1)
	logs.MyErrorLog.CheckFatallnError("约租设置失败:", err)

	// 注册服务节点
	ipPort := constants.Configs["serverNode.ip"] + ":" + constants.Configs["serverNode.port"]
	nodeKey := "servers||" + ipPort
	nodeInfo := "IP:" + constants.Configs["serverNode.ip"] + ",PORT:" + constants.Configs["serverNode.port"] + ";setTime:" + time.Now().Format("2006-01-02 15:04:05.9999")
	_, err = cli.Put(context.Background(), nodeKey, nodeInfo, clientv3.WithLease(leaseGrantResp.ID))
	logs.MyInfoLog.CheckFatallnError("注册服务节点:"+ipPort+"失败:", err)
	logs.MyInfoLog.Println("注册服务节点:" + nodeKey)

	// 注册服务接口
	for _, service := range service.Services {
		sKey := nodeKey + "||" + service.URL + "||" + service.MethodType
		sValue := "service_id:" + service.ServiceId + ";service_name:" + service.ServiceName + ";method:" + service.Method

		_, err = cli.Put(context.TODO(), sKey, sValue, clientv3.WithLease(leaseGrantResp.ID))
		logs.MyErrorLog.CheckFatallnError("注册服务接口:"+sKey+"失败:", err)
		logs.MyInfoLog.Println("注册服务接口:" + sKey)
	}

	logs.MyInfoLog.Println("注册节点信息到Etcd完成!")

	registerChan <- true

	// 保持连接
	for {
		cli.KeepAliveOnce(context.Background(), leaseGrantResp.ID)
		//		logs.MyDebugLog.Println("连接到etcd正常...")
		time.Sleep(100 * time.Millisecond)
	}
}

// 注册监听RPC
func registerRpc() {
	logs.MyDebugLog.Println("register rpc now...")

	serviceNode := new(service.ServiceNode)
	err := rpc.Register(serviceNode)
	logs.MyErrorLog.CheckFatallnError("", err)

	rpc.HandleHTTP()
	err = http.ListenAndServe(":"+constants.Configs["serverNode.port"], nil)
	logs.MyErrorLog.CheckFatallnError("", err)
	logs.MyDebugLog.Println("register rpc success...")
}

// 服务中心的监控Etcd代码
func RegisterToEtcd1() {
	go func() {
		logs.MyInfoLog.Println("开始监听Etcd...")

		// 连接到etcd
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{constants.Configs["etcd.url"]},
			DialTimeout: 100 * time.Second,
		})
		logs.MyErrorLog.CheckPaniclnError("监听程序创建etcd clientv3失败:", err)
		defer cli.Close()

		for {
			wChan := cli.Watch(context.TODO(), "servers", clientv3.WithPrefix())
			wc := <-wChan
			fmt.Printf("监听到结果,canceled:%t,created:%t,Header:%d,isProgressNotify:%t \n", wc.Canceled, wc.Created, wc.Header, wc.IsProgressNotify())
			for _, e := range wc.Events {
				fmt.Printf("isCreate:%t,isModify:%t,Type:%s,Key:%s,Value:%s \n", e.IsCreate(), e.IsModify(), e.Type, e.Kv.Key, e.Kv.Value)
			}
		}

		logs.MyInfoLog.Println("监听Etcd结束!")
	}()
}
