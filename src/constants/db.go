package constants

import (
	"database/sql"
	"fmt"
	"logs"

	_ "github.com/go-sql-driver/mysql"
)

// Postgre database
var PostgreDB sql.DB

// MySQL database
var MySQLDB *sql.DB

var MyName string = "Huanggh sb"

// 初始化
func init() {
	initMySQLDB()
}

// 初始化数据库
func initMySQLDB() {
	MySQLDB, err := sql.Open(Configs["database"], Configs["mysql.url"])
	logs.MyErrorLog.CheckFatallnError("打开MySQL失败:", err)

	fmt.Println("数据库初始化成功!", MySQLDB)

	err = MySQLDB.Ping()
	if err != nil {
		fmt.Println("数据库初始化未成功", err)
	}

	fmt.Println("status:", MySQLDB.Stats().OpenConnections)

	// do something
	rows, err := MySQLDB.Query("select id,name,time from sys_user where name like ?", "%huanggh%")
	logs.MyInfoLog.CheckPrintlnError("query:", err)

	for rows.Next() {
		var id int
		var name string
		var date string
		rows.Scan(&id, &name, &date)

		fmt.Printf("id:%d,name:%s,Time:%s \n", id, name, date)
	}

}
