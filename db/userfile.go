package db

import (
	"fmt"
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
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`," +
			"`upload_at`)values(?,?,?,?,?)" )
	defer stmt.Close()
	if err != nil{
		return false
	}

	_, err = stmt.Exec(username,filehash,filename,filesize,time.Now())
	if err != nil{
		return false
	}
	return true
}
//QueryUserFileMeta 批量获取为用户文件信息
func QueryUserFileMeta(username string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err := rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}
