package main

import (
	"log"
	"his/handle"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handle.UploadHandler)
	http.HandleFunc("/file/upload/succ", handle.UploadSuccHandler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Panic(err)
	}
}
