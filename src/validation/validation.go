// 接口入参校验
package validation

import (
	"logs"
)

// 接口入参校验
func CheckInterfaceInParams(methodName string, inParams map[string]interface{}) error {
	logs.MyDebugLog.Println("接口[", methodName, "]入参校验开始.")

	// doing checking...
	logs.MyDebugLog.Println("校验入参:", inParams)

	logs.MyDebugLog.Println("接口[", methodName, "]入参校验结束!")
	return nil
}
