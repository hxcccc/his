package mysql

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
func init() {
	db, _ = sql.Open("mysql", "root:root@tcp(10.8.149.172:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(100)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql,err:", err.Error())
		os.Exit(1)
	}
}
//DBConn 返回数据库连接
func DBConn() *sql.DB {
	return db
}
