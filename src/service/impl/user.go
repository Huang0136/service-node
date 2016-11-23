package impl

import (
	"beans"
	"constants"
	"crypto/sha512"
	"fmt"
	"io"
	"logs"
	"strconv"
	"time"
	"utils"
)

// 登录
func (si *ServiceImpl) Login() (msg string, err error) {
	userName := si.InParams["USER_NAME"].(string) // 输入的登录名
	password := si.InParams["PASSWORD"].(string)  // 输入的密码

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

	fmt.Printf("结果,id:%d,un:%s,pw:%s,time:%s,status:%s\n", id, un, dbPassword, ts, status)

	if id == 0 && dbPassword == "" {
		msg = "用户名/密码不正确"
		return
	}

	// 密码加密存储，sha512(sha512(password),salt)
	shaHash1 := sha512.New()
	io.WriteString(shaHash1, password)

	str1 := fmt.Sprintf("%x", shaHash1.Sum(nil)) // 第一次加密

	shaHash2 := sha512.New()
	io.WriteString(shaHash2, str1)
	io.WriteString(shaHash2, strconv.Itoa(id))

	str2 := fmt.Sprintf("%x", shaHash2.Sum(nil)) // 第二次加密

	fmt.Printf("用户登录,username:%s,password:%s,sha512加密后:%s \n", userName, password, str2)

	if str2 == dbPassword {
		token := utils.CreateToken()
		msg = "{\"message\":\"登录成功\",\"TOKEN\":\"" + token + "\"}"
		return
	} else {
		msg = "用户名/密码不正确"
		return
	}
}

// 根据用户Id获取用户
func (si *ServiceImpl) GetUserByUserId() (u beans.User, other map[string]interface{}, err error) {
	userId := si.InParams["USER_ID"].(string)
	fmt.Println("业务方法接收到参数:", si.InParams)

	t1 := time.Now().UnixNano()
	sql := "select id,user_name,time from sys_user where id = ?"

	stat, err := constants.MySQLDB.Prepare(sql)
	logs.MyErrorLog.CheckPrintlnError("prepare:", err)

	rs, err := stat.Query(userId)
	logs.MyErrorLog.CheckPrintlnError("query:", err)

	//	var b bytes.Buffer

	var list []beans.User
	for rs.Next() {
		var id int
		var name string
		var time1 []uint8

		err = rs.Scan(&id, &name, &time1)
		logs.MyInfoLog.CheckPrintlnError("scan value:", err)

		t1, _ := time.Parse("2006-01-02 15:04:05.999999999", string(time1)) //2006-01-02 15:04:05.99999999

		u := beans.User{
			Id: id, Name: name, Time: t1,
		}

		list = append(list, u)
	}

	/*
		b.WriteString("[")
		ll := len(list)
		for i, v := range list {
			b.WriteString("{")
			b.WriteString("\"id\":\"")
			b.WriteString(strconv.Itoa(v.Id))
			b.WriteString("\",\"user_name\":\"")
			b.WriteString(v.Name)
			b.WriteString("\",\"time\":\"")
			b.WriteString(v.Time.Format("2006-01-02 15:04:05.9999"))
			b.WriteString("\"")
			if ll == i+1 {
				b.WriteString("}")
			} else {
				b.WriteString("},")
			}

		}
		b.WriteString("]")
	*/
	t2 := time.Now().UnixNano()

	u = list[0]
	other = make(map[string]interface{})
	other["SQL_EXCU_TIME_BEGIN"] = t1
	other["SQL_EXCU_TIME_END"] = t2
	return

}

// 添加User
func (si *ServiceImpl) Add(user beans.User) string {
	sql := "insert into user(id,user_name,password,time,status)values(?,?,?,?,?)"
	stat, err := constants.MySQLDB.Prepare(sql)
	logs.MyInfoLog.CheckPrintlnError("", err)

	rs, err := stat.Exec(user.Id, user.Name, user.Time)
	logs.MyInfoLog.CheckPrintlnError("", err)

	row, err := rs.RowsAffected()
	logs.MyInfoLog.CheckPrintlnError("", err)

	return "add success:" + strconv.Itoa(int(row))
}
