package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logs"
	"os"
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
	serviceConfig, err := os.Open("./config/service.json")
	logs.MyDebugLog.CheckFatallnError("read service config file error:", err)

	sByte, err := ioutil.ReadAll(serviceConfig)
	logs.MyDebugLog.CheckFatallnError("read service config to byte error:", err)

	err = json.Unmarshal(sByte, &Services)
	logs.MyDebugLog.CheckFatallnError("convert byte to json error:", err)

	// 打印服务接口
	PrintlnServices()
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
