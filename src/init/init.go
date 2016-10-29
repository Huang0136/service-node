// 项目启动时初始化
package init

import (
	"encoding/json"
	"io/ioutil"
	"logs"
	"os"
	"service"
)

func init() {
	serviceConfig, err := os.Open("./config/service.json")
	logs.MyDebugLog.CheckFatallnError("read service config file error:", err)

	sByte, err := ioutil.ReadAll(serviceConfig)
	logs.MyDebugLog.CheckFatallnError("read service config to byte error:", err)

	err = json.Unmarshal(sByte, &service.Services)
	logs.MyDebugLog.CheckFatallnError("convert byte to json error:", err)
}
