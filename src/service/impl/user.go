package impl

import (
	"constants"
	"logs"
)

// 根据用户Id获取用户
func (si *ServiceImpl) GetUserByUserId(userId string) {
	sql := "select * from user where id = ?"

	stat, err := constants.MySQLDB.Prepare(sql)
	logs.MyErrorLog.CheckPrintlnError(err)

	rs, err := stat.Exec(userId)
	logs.MyErrorLog.CheckPrintlnError(err)
	
	rs.

}

func (si *ServiceImpl) TT(ui string) string {

	return ""
}
