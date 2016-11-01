// 接口入参校验
package validation

import (
	"logs"
)

//
func CheckInterfaceInParams(methodName string) error {
	logs.MyDebugLog.Println("接口【", methodName, "】入参校验开始...")

	// doing checking...

	logs.MyDebugLog.Println("接口【", methodName, "】入参校验结束!")
	return nil
}
