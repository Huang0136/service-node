package impl

import (
	"beans"
	"constants"
	"fmt"
	"logs"
	"strconv"
	"time"
	"utils"
)

// 用户登录
//
func (si *ServiceImpl) Login() (result map[string]interface{}, other map[string]interface{}, err error) {
	// 入参
	userName := si.InParams["USER_NAME"].(string) // 输入的登录名
	password := si.InParams["PASSWORD"].(string)  // 输入的密码

	tQuery1 := time.Now()
	sql := "select t.id,t.user_name,t.password,t.time,t.status from sys_user t where t.user_name = ? "
	stat, err := constants.MySQLDB.Prepare(sql)
	if err != nil {
		logs.MyInfoLog.Panicln("", err)
		return
	}

	row := stat.QueryRow(userName)

	var id int
	var dbPassword string
	var un string
	var ts []uint8
	var status string
	err = row.Scan(&id, &un, &dbPassword, &ts, &status)
	if err != nil {
		logs.MyInfoLog.Panicln("", err)
		return
	}
	tQuery2 := time.Now()

	fmt.Printf("结果,id:%d,un:%s,pw:%s,time:%s,status:%s\n", id, un, dbPassword, ts, status)

	if id == 0 && un == "" {
		result["message"] = "用户名/密码不正确"
		return
	}

	// 密码加密
	tPwEncrypt1 := time.Now()
	str2 := utils.EncryptPassword(password, strconv.Itoa(id))
	fmt.Printf("用户登录,username:%s,password:%s,sha512加密后:%s \n", userName, password, str2)
	tPwEncrypt2 := time.Now()

	// 返回结果(性能分析)
	other = make(map[string]interface{})
	other["TIME_SQL_SELECT_USER"] = tQuery2.Sub(tQuery1).Seconds()
	other["TIME_ENCRYPT_PASSWORD"] = tPwEncrypt2.Sub(tPwEncrypt1).Seconds()

	// 返回结果
	result = make(map[string]interface{})
	if str2 == dbPassword {
		result["message"] = "登录成功"
		result["token"] = utils.CreateToken()
		return
	} else {
		result["message"] = "用户名/密码不正确"
		return
	}
}

// 根据用户Id获取用户
func (si *ServiceImpl) GetUserByUserId() (u beans.User, other map[string]interface{}, err error) {
	userId := si.InParams["USER_ID"].(string)

	tQuery1 := time.Now()
	sql := "select id,user_name,time from sys_user where id = ?"

	stat, err := constants.MySQLDB.Prepare(sql)
	logs.MyErrorLog.CheckPrintlnError("prepare:", err)

	rs, err := stat.Query(userId)
	logs.MyErrorLog.CheckPrintlnError("query:", err)

	var list []beans.User
	for rs.Next() {
		var id int
		var name string
		var time1 []uint8

		err = rs.Scan(&id, &name, &time1)
		logs.MyInfoLog.CheckPrintlnError("scan value:", err)

		t1, _ := time.Parse("2006-01-02 15:04:05.999999999", string(time1)) //2006-01-02 15:04:05.99999999

		uTemp := beans.User{
			Id: id, Name: name, Time: t1,
		}

		list = append(list, uTemp)
	}
	tQuery2 := time.Now()

	// 返回结果(性能分析)
	other = make(map[string]interface{})
	other["TIME_SQL_SELECT_USER"] = tQuery2.Sub(tQuery1).Seconds()

	u = list[0]
	return
}

// 新增用户
func (si *ServiceImpl) Add() (result map[string]interface{}, other map[string]interface{}, err error) {
	// 入参
	id := si.InParams["user_id"].(string)
	userName := si.InParams["user_name"].(string)
	pw := si.InParams["password"].(string)
	status := 0

	tEncryptPw1 := time.Now()
	pwEncrypt := utils.EncryptPassword(pw, id)
	tEncryptPw2 := time.Now()

	tQuery1 := time.Now()
	sql := "insert into user(id,user_name,password,time,status)values(?,?,?,?,?)"
	stat, err := constants.MySQLDB.Prepare(sql)
	logs.MyInfoLog.CheckPrintlnError("", err)

	rs, err := stat.Exec(id, userName, pwEncrypt, time.Now(), status)
	logs.MyInfoLog.CheckPrintlnError("", err)

	row, err := rs.RowsAffected()
	logs.MyInfoLog.CheckPrintlnError("", err)

	tQuery2 := time.Now()

	// 返回结果(性能分析)
	other = make(map[string]interface{})
	other["TIME_SQL_SELECT_USER"] = tQuery2.Sub(tQuery1).Seconds()
	other["TIME_ENCRYPT_PASSWORD"] = tEncryptPw2.Sub(tEncryptPw1).Seconds()

	// 返回结果
	result = make(map[string]interface{})
	result["MESSAGE"] = "新增用户成功"
	result["Affected"] = row

	return
}
