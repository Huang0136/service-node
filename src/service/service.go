package service

import (
	"bytes"
	"constants"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"logs"
	"net/http"
	"net/rpc"
	"os"
	"reflect"
	"service/impl"
	"time"
	"validation"
)

// 服务节点
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
	Regexp      bool       `json:"regexp"`
	MethodType  string     `json:"method_type"`
	Method      string     `json:"method"`
	ContentType string     `json:"content_type"`
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
	// 读取服务接口配置文件
	readServiceConfig()

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

// 统一的rpc处理方法
func (serverNode *ServiceNode) RpcCallHandler(req Req, resp *Resp) error {
	tNodeBegin := time.Now()                            // 开始处理时间
	methodName, _ := req.Params["METHOD_NAME"].(string) // 方法名称

	// 业务服务实现
	si := new(impl.ServiceImpl)
	si.InParams = req.Params
	fmt.Println("请求参数:", si.InParams)

	// 接口入参校验
	err := validation.CheckInterfaceInParams(methodName, si.InParams)
	if err != nil {
		return err
	}

	// 封装入参
	tBusineBegin := time.Now() //反射调用开始
	var inReflectValues []reflect.Value = make([]reflect.Value, 0)

	// 反射调用
	funcTemp := reflect.ValueOf(si).MethodByName(methodName)
	rfValues := funcTemp.Call(inReflectValues)
	tBusineEnd := time.Now() // 反射调用结束

	fmt.Println("业务执行结果:", rfValues)
	fmt.Println("返回参数个数:", len(rfValues))
	for i, rfv := range rfValues {
		fmt.Printf("第%d个参数,类型:%s,值:%v\n", i, rfv.Kind().String(), rfv)
		//		kt, _ := strconv.Atoi(rfv.Kind().String())
		//		kindType := uint(kt)
		switch rfv.Kind() {
		case reflect.Interface:
			fmt.Printf("返回值类型:interface,%v\n", rfv)
			break
		case reflect.Struct:
			fmt.Printf("返回值类型:Struct,%v\n", rfv)
			break
		case reflect.Map:
			fmt.Printf("返回值类型:Map,%v\n", rfv)
			break
		case reflect.String:
			fmt.Printf("返回值类型:String,%v\n", rfv)
			break
		case reflect.Slice:
			fmt.Printf("返回值类型:Slice,%v\n", rfv)
			break
		case reflect.Array:
			fmt.Printf("返回值类型:Array,%v\n", rfv)
			break
		case reflect.Int:
			fmt.Printf("返回值类型:Int,%v\n", rfv)
			break
		default:
			fmt.Printf("未知的返回值类型:%s\n", rfv.Kind())
		}
	}

	/*
		Invalid Kind = iota
		Bool
		Int
		Int8
		Int16
		Int32
		Int64
		Uint
		Uint8
		Uint16
		Uint32
		Uint64
		Uintptr
		Float32
		Float64
		Complex64
		Complex128
		Array
		Chan
		Func
		Interface
		Map
		Ptr
		Slice
		String
		Struct
		UnsafePointer
	*/

	// 返回结果
	resp.Params = make(map[string]interface{})
	//resp.Params["result"] = rfValues[0] //.String() []User、User、string、int、interface、map[string]interface

	rsType := rfValues[0].Kind()
	if rsType == reflect.String {
		resp.Params["result"] = rfValues[0].String()
	} else if rsType == reflect.Slice {
		resp.Params["result"] = rfValues[0].Bytes()
	} else if rsType == reflect.Struct {

		resp.Params["result"] = rfValues[0].String()
	}

	tNodeEnd := time.Now()

	resp.Params["SERVER_NODE"] = constants.Configs["serverNode.ip"] + ":" + constants.Configs["serverNode.port"]
	resp.Params["NODE_BEGIN_TIME"] = tNodeBegin.UnixNano()
	resp.Params["NODE_END_TIME"] = tNodeEnd.UnixNano()
	resp.Params["BUSINE_BEGIN_TIME"] = tBusineBegin.UnixNano()
	resp.Params["BUSINE_END_TIME"] = tBusineEnd.UnixNano()

	// 业务执行时间
	log.Printf("SQL执行时间:%f,节点执行时间:%f\n", tBusineEnd.Sub(tBusineBegin).Seconds(), tNodeEnd.Sub(tNodeBegin).Seconds())
	return nil
}

// 注册rpc
func RegisterRpc() {
	logs.MyDebugLog.Println("register rpc...")

	serviceNode := new(ServiceNode)
	err := rpc.Register(serviceNode)
	logs.MyErrorLog.CheckFatallnError("", err)

	rpc.HandleHTTP()
	err = http.ListenAndServe(":"+constants.Configs["serverNode.port"], nil)
	logs.MyErrorLog.CheckFatallnError("", err)
}

// 读取服务接口配置文件
func readServiceConfig() {
	serviceConfig, err := os.Open("./config/service.json")
	logs.MyDebugLog.CheckFatallnError("read service config file error:", err)

	sByte, err := ioutil.ReadAll(serviceConfig)
	logs.MyDebugLog.CheckFatallnError("read service config to byte error:", err)

	err = json.Unmarshal(sByte, &Services)
	logs.MyDebugLog.CheckFatallnError("convert byte to json error:", err)
}
