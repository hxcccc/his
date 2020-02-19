package db

import (
	"his/db/mysql"
	"time"
)

//UserFile:用户文件表结构体
type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}
//OnUserFileUploadFinished: 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user_file (`username`,`file_sha1`,`filename`,`filesize`," +
			"`update_at`)values(?,?,?,?,?)" )
	if err != nil{
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username,filehash,filename,filesize,time.Now())
	if err != nil{
		return false
	}
	return true
}
