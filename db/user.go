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
//UserSignIn 判断密码是否一致
func UserSignIn(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return  false
	}else if rows == nil {
		fmt.Printf("username:%s not found\n", username)
		return false
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}

	return false
}
//UpdateToken 刷新token
func UpdateToken(username string, token string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?.?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}
