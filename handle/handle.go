package handle

import (
	"log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		//返回上传页面
		data, err := ioutil.ReadFile("./static/view/index.html")
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

		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			log.Panic(err)
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			log.Panic(err)
		}

		http.Redirect(w, r, "/file/upload/succ", http.StatusFound)
	}
}
//UploadSuccHandler: 上传已完成
func UploadSuccHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload finished")
}
