package handle

import (
	"encoding/json"
	"fmt"
	"his/db"
	"his/meta"
	"his/util"
	"log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UserHome(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./static/view/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(data))
}

func UploadHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		//返回上传页面
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件流及存储到本地
		file, head, err := r.FormFile("file")
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName:head.Filename,
			Location:"/tmp/" + head.Filename,
			UploadAt:time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			log.Panic(err)
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			log.Panic(err)
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		_ = meta.UpdateFileMetaDB(fileMeta)

		//TODO: 更新用户文件表记录
		r.ParseForm()
		username := r.Form.Get("username")
		suc := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName,
			fileMeta.FileSize)
		if suc{
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		}else {
			w.Write([]byte("Upload Failed"))
		}

		//http.Redirect(w, r, "/file/upload/succ", http.StatusFound)
	}
}
//UploadSuccHandler: 上传已完成
func UploadSuccHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload finished")
}
//GetFileMetaHandler 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form["filehash"][0]
	fMeta,err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w,"no such file")
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
//DownLoadHandler 文件下载
func DownLoadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm, err := meta.GetFileMetaDB(fsha1)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "no such file")
		return
	}
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/otect-stream")
	w.Header().Set("Content-Descrption", "attachment;filename=\"" + fm.FileName+"\"")
	w.Write(data)
}
//FileMetaUploadHandler 更新元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opTpey := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opTpey != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	curFileMeta.FileName = newFileName
	os.Rename(curFileMeta.Location, "/tmp/" + newFileName)
	curFileMeta.Location = "/tmp/" + newFileName
	curFileMeta.UploadAt = time.Now().Format("2006-1-2 15:04:05")
	_ = meta.ReplaceFileMetaDB(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
//FileDelHandler 删除文件
func FileDelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")

	fm := meta.GetFileMeta(fileSha1)
	os.Remove(fm.Location)

	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "delete success")
}
