package service

import (
	"bytes"
	"constants"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"logs"
	"net/http"
	"net/rpc"
	"os"
	"reflect"
	"service/impl"
	"strings"
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

func init() {
	// 读取服务接口配置文件
	readServiceConfig()

}

// 统一的rpc处理方法
func (serverNode *ServiceNode) RpcCallHandler(req Req, resp *Resp) error {
	tNode1 := time.Now()                                // 节点处理开始时间
	methodName, _ := req.Params["METHOD_NAME"].(string) // 方法名称

	// 业务服务实现
	si := new(impl.ServiceImpl)
	si.InParams = req.Params

	// 服务配置
	service := findService(methodName)

	// 接口入参校验
	tParamCheck1 := time.Now()
	err := validation.CheckInterfaceInParams(methodName, si.InParams)
	if err != nil {
		return err
	}
	tParamCheck2 := time.Now()

	// 封装入参
	tReflectCall1 := time.Now() //反射调用开始时间
	var inReflectValues []reflect.Value = make([]reflect.Value, 0)

	// TODO 封装给方法的入参

	// 反射调用
	funcTemp := reflect.ValueOf(si).MethodByName(methodName)
	rfValues := funcTemp.Call(inReflectValues)
	tReflectCall2 := time.Now() // 反射调用结束时间

	fmt.Println("业务执行结果:", rfValues)

	// 返回结果（给服务中心）
	resp.Params = make(map[string]interface{})

	// 结果:业务执行是否有错
	rsError := rfValues[2]
	if rsError.Interface() != nil {
		errTemp := rsError.Interface().(error)
		err1 := errors.New("业务执行失败," + errTemp.Error())
		log.Panicln(err1)
		return err1
	}

	// 结果:业务返回
	rsReturn := rfValues[0]

	switch service.ContentType {
	case "json":
		bJson, err := json.Marshal(rsReturn.Interface())
		if err != nil {
			err1 := errors.New("转换JSON失败," + err.Error())
			log.Panicln(err1)
			return err1
		}

		resp.Params["RESULT"] = string(bJson)
		break
	case "file":
		resp.Params["RESULT"] = rsReturn.Bytes()
		break
	case "static":
		resp.Params["RESULT"] = rsReturn.Bytes()
		break
	default:
		log.Panicln("未知的返回类型")
	}

	// 结果:用于性能调优
	rsDefault := rfValues[1]

	var str1 bytes.Buffer
	str1.WriteString("业务执行耗时,")
	for _, m := range rsDefault.MapKeys() {
		k1 := m.String()
		v1 := rsDefault.MapIndex(m).Interface()
		strV1 := fmt.Sprint(v1.(float64))

		str1.WriteString(k1)
		str1.WriteString(":")
		str1.WriteString(strV1)
		str1.WriteString(",")

		resp.Params["HANDLER_"+k1] = strV1 // 业务处理各耗时情况
	}

	str2 := strings.TrimRight(str1.String(), ",")
	fmt.Println(str2)

	resp.Params["HANDLER_SERVER_NODE"] = constants.Configs["serverNode.ip"] + ":" + constants.Configs["serverNode.port"]
	resp.Params["TIME_SERVICE_PARAMS_CHECK"] = fmt.Sprint(tParamCheck2.Sub(tParamCheck1).Seconds()) // 入参校验耗时
	resp.Params["TIME_REFLECT_CALL"] = fmt.Sprint(tReflectCall2.Sub(tReflectCall1).Seconds())       // 反射调用耗时

	tNode2 := time.Now()                                                      // 节点处理结束时间
	resp.Params["TIME_NODE_TOTAL"] = fmt.Sprint(tNode2.Sub(tNode1).Seconds()) // 节点处理总耗时

	// 耗时
	fmt.Printf("节点处理总耗时:%f,反射调用耗时:%f,入参校验耗时:%f\n", tNode2.Sub(tNode1).Seconds(), tReflectCall2.Sub(tReflectCall1).Seconds(), tParamCheck2.Sub(tParamCheck1).Seconds())
	return nil
}

// 根据method查找接口配置信息
func findService(methodName string) (s Service) {
	for _, sTemp := range Services {
		if methodName == sTemp.Method {
			s = sTemp
			break
		}
	}
	return
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
