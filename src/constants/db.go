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

func init() {
	initMySQLDB()
}

func initMySQLDB() {

	MySQLDB, err := sql.Open("mysql", "root:Huang0136@tcp/server-node?charset=utf-8")
	logs.MyErrorLog.CheckFatallnError("", err)

	fmt.Println("数据库初始化成功!", MySQLDB)

}
