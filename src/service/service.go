package service

import (
	"bytes"
	"constants"
	"fmt"
	"logs"
	"net/http"
	"net/rpc"
	"reflect"
	"service/impl"
	"validation"
)

// 服务
type ServiceNode int

// 请求参数
type Req struct {
	Params map[string]interface{}
}

// 返回数据
type Resp struct {
	Params map[string]interface{}
}

// 服务接口结构体
type Service struct {
	ServiceId   string     `json:"service_id"`
	ServiceName string     `json:"service_name"`
	URL         string     `json:"url"`
	MethodType  string     `json:"method_type"`
	Method      string     `json:"method"`
	Remark      string     `json:"remark"`
	InParams    []InParam  `json:"in_params"`
	OutParams   []OutParam `json:"out_params"`
}

// 入参结构体
type InParam struct {
	ParamCode string `json:"param_code"`
	ParamName string `json:"param_name"`
	ParamType string `json:"param_type"`
	Require   bool   `json:"requrie"`
	Length    int    `json:"length"`
	Remark    string `json:"remark"`
}

// 出参结构体
type OutParam struct {
	ParamCode string `json:"param_code"`
	ParamName string `json:"param_name"`
	ParamType string `json:"param_type"`
	Remark    string `json:"remark"`
}

// 服务接口列表
var Services []Service = make([]Service, 0)

//
func init() {

	// 打印服务接口
	//	PrintlnServices()
}

// 打印服务接口列表
func PrintlnServices() {
	for _, s := range Services {
		fmt.Println(ServiceToStr(s))
	}
}

// 将服务转换成字符串
func ServiceToStr(s Service) string {
	b := bytes.Buffer{}

	b.WriteString("接口:" + s.ServiceName + ",")
	b.WriteString("方法:" + s.Method)

	return b.String()
}

// 统一的rpc调用处理方法
func (serverNode *ServiceNode) RpcCallHandler(req Req, resp *Resp) error {
	// 方法名称
	methodName, _ := req.Params["METHOD_NAME"].(string)

	fmt.Println("MethodName:", methodName)

	si := new(impl.ServiceImpl)
	si.InParams = req.Params

	funcTemp := reflect.ValueOf(si).MethodByName(methodName)

	// 封装入参
	//	var inValues []reflect.Value

	// 接口入参校验
	err := validation.CheckInterfaceInParams(methodName)
	if err != nil {
		return err
	}

	// 反射调用
	rfValues := funcTemp.Call(nil)

	// 封装出参
	resp.Params["result"] = string(rfValues[0].Bytes())
	return nil
}

// 注册rpc
func RegisterRpc() {
	logs.MyDebugLog.Println("register rpc now...")

	serviceNode := new(ServiceNode)
	err := rpc.Register(serviceNode)
	logs.MyErrorLog.CheckFatallnError("", err)

	rpc.HandleHTTP()
	err = http.ListenAndServe(":"+constants.Configs["serverNode.port"], nil)
	logs.MyErrorLog.CheckFatallnError("", err)
	logs.MyDebugLog.Println("register rpc success...")
}
