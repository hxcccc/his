package main

import (
	"log"
	"his/handle"
	"net/http"
)

func main() {
	//文件存取接口
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/upload", handle.HTTPinterceptor(handle.UploadHandler))
	http.HandleFunc("/file/upload/succ", handle.HTTPinterceptor(handle.UploadSuccHandler))
	http.HandleFunc("/file/meta", handle.HTTPinterceptor(handle.GetFileMetaHandler))
	http.HandleFunc("/file/download", handle.HTTPinterceptor(handle.DownLoadHandler))
	http.HandleFunc("/file/update", handle.HTTPinterceptor(handle.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete", handle.HTTPinterceptor(handle.FileDelHandler))
	http.HandleFunc("/file/query", handle.HTTPinterceptor(handle.FileQueryHandler))
	http.HandleFunc("/user/home", handle.UserHome)
	http.HandleFunc("/file/fastupload", handle.HTTPinterceptor(handle.TryFastUploadHandler))


	//用户相关接口
	http.HandleFunc("/user/signup", handle.SignUpHandler)
	http.HandleFunc("/user/signin", handle.SignInHandler)
	http.HandleFunc("/user/info", handle.HTTPinterceptor(handle.UserInfoHandler))

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Panic(err)
	}
}
