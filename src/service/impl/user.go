package impl

import (
	"bytes"
	"constants"
	"logs"
	"strconv"
	"time"
)

type User struct {
	Id   int
	Name string
	Time time.Time
}

// 根据用户Id获取用户
func (si *ServiceImpl) GetUserByUserId() (msg string, err error) {
	userId := si.InParams["USER_ID"].(string)

	sql := "select * from sys_user where id = ?"

	stat, err := constants.MySQLDB.Prepare(sql)
	logs.MyErrorLog.CheckPrintlnError("prepare:", err)

	rs, err := stat.Query(userId)
	logs.MyErrorLog.CheckPrintlnError("query:", err)

	var b bytes.Buffer

	var list []User
	for rs.Next() {
		var id int
		var name string
		var time1 []uint8

		err = rs.Scan(&id, &name, &time1)
		logs.MyInfoLog.CheckPrintlnError("scan value:", err)

		t1, _ := time.Parse("2006-01-02 15:04:05.999999999", string(time1)) //2006-01-02 15:04:05.99999999

		u := User{
			Id: id, Name: name, Time: t1,
		}

		list = append(list, u)
	}

	b.WriteString("[")
	ll := len(list)
	for i, v := range list {
		b.WriteString("{")
		b.WriteString("\"id\":\"")
		b.WriteString(strconv.Itoa(v.Id))
		b.WriteString("\",\"name\":\"")
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

	return b.String(), nil

}

// 添加User
func (si *ServiceImpl) Add(user User) string {
	sql := "insert into user(id,name,time)values(?,?,?)"
	stat, err := constants.MySQLDB.Prepare(sql)
	logs.MyInfoLog.CheckPrintlnError("", err)

	rs, err := stat.Exec(user.Id, user.Name, user.Time)
	logs.MyInfoLog.CheckPrintlnError("", err)

	row, err := rs.RowsAffected()
	logs.MyInfoLog.CheckPrintlnError("", err)

	return "add success:" + strconv.Itoa(int(row))
}
