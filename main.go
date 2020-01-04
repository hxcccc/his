package main

import (
	"log"
	"his/handle"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/upload", handle.UploadHandler)
	http.HandleFunc("/file/upload/succ", handle.UploadSuccHandler)
	http.HandleFunc("/file/meta", handle.GetFileMetaHandler)
	http.HandleFunc("/file/download", handle.DownLoadHandler)
	http.HandleFunc("/file/update", handle.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handle.FileDelHandler)
	http.HandleFunc("/user/signup", handle.SignUpHandler)
	http.HandleFunc("/user/signin", handle.SignInHandler)
	http.HandleFunc("/user/home", handle.UserHome)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Panic(err)
	}
}
