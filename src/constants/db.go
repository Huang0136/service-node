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
var MySQLDB sql.DB

// 初始化
func init() {
	initMySQLDB()
}

// 初始化数据库
func initMySQLDB() {
	MySQLDB, err := sql.Open(Configs["database"], Configs["mysql.url"])
	logs.MyErrorLog.CheckFatallnError("", err)

	fmt.Println("数据库初始化成功!", MySQLDB)
}
