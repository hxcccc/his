package db

import (
	"fmt"
	"his/db/mysql"
)
//UserSignUp 通过用户名及密码完成用户表注册
func UserSignUp(username string, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err:", err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:", err.Error())
		return false
	}
	if rf, err := ret.RowsAffected();err ==nil && rf >0 {
		return true
	}
	return false
}
