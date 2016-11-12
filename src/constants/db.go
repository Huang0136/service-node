package constants

import (
	"database/sql"
	"logs"

	_ "github.com/go-sql-driver/mysql"
)

// Postgre database
var PostgreDB *sql.DB

// MySQL database
var MySQLDB *sql.DB

// 初始化
func init() {
	initMySQLDB()
	doMySQLExec()
}

// 初始化数据库
func initMySQLDB() {
	mysqlDB, err := sql.Open(Configs["database"], Configs["mysql.url"])
	logs.MyErrorLog.CheckFatallnError("打开MySQL失败:", err)

	logs.MyInfoLog.Println("数据库初始化成功!", mysqlDB)
	err = mysqlDB.Ping()
	if err != nil {
		logs.MyInfoLog.Println("数据库初始化未成功", err)
	}

	MySQLDB = mysqlDB
}

// Test: select 1
func doMySQLExec() {
	rows, err := MySQLDB.Query("select 1 ")
	logs.MyInfoLog.CheckPrintlnError("query: select 1", err)

	for rows.Next() {
		var id int
		rows.Scan(&id)

		logs.MyInfoLog.Println("[sql:select 1] result:", id)
	}
}
