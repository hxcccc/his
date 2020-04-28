package handle

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	muredis "his/cache/redis"
	"his/db"
	"his/util"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)
//MultipartUploadInfo 初始化信息
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}

//初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}
	//获得redis链接
	rConn := muredis.RedisPool().Get()
	defer rConn.Close()
	//生成分开上传的初始化信息
	upinfo := MultipartUploadInfo{
		FileHash:filehash,
		FileSize:filesize,
		UploadID:username+fmt.Sprint("%x", time.Now().UnixNano()),
		ChunkSize:5*1024*1024,
		ChunkCount:int(math.Ceil(float64(filesize/(5*1024*1024)))),
	}
	//初始化信息写入redis缓存
	rConn.Do("HSET", "MP_"+upinfo.UploadID, "chunkcount", upinfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upinfo.UploadID, "filehash", upinfo.FileHash)
	rConn.Do("HSET", "MP_"+upinfo.UploadID, "filesize", upinfo.FileSize)
	//将初始化信息返回给客户端
	w.Write(util.NewRespMsg(0, "OK", upinfo).JSONBytes())
}

// UploadPartHandler  上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//解析用户请求参数
	r.ParseForm()
	//	username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 获得redis连接池中的一个连接
	rConn := muredis.RedisPool().Get()
	defer rConn.Close()

	//获得文件句柄，用于存储分块内容
	fpath := "/tmp/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	//更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

//CompleteUploadHandler 通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	//解析请求参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")
	//获取redis连接池的一个连接
	rConn := muredis.RedisPool().Get()
	defer rConn.Close()
	//通过uploadid查询redis，判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		}else if strings.HasPrefix(k, "chkidx_") && v == "1"{
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}
	//合并分块
	//后续补充.....
	//更新唯一文件表以及用户文件表
	fsize, _ := strconv.ParseInt(filesize, 10, 64)
	db.OnfileUploadFinished(filehash, filename, fsize,"")
	db.OnUserFileUploadFinished(username, filehash, filename, fsize)
	//响应处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
