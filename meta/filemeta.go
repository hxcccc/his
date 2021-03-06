package meta

import "his/db"

//FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}
//UpdateFileMeta 新增或者更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnfileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

func ReplaceFileMetaDB(meta FileMeta) bool {
	return db.UpdateFileDB(meta.FileSha1, meta.FileName,meta.Location)
}

//GetFileMeta 通过sha1值获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func GetFileMetaDB(fileSha1 string) (FileMeta,error)  {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{},err
	}
	fmeta := FileMeta{
		FileSha1:tfile.FileHash,
		FileSize:tfile.FileSize.Int64,
		FileName:tfile.FileName.String,
		Location:tfile.FileAddr.String,
		UploadAt:string(tfile.FileUpdataAt),
	}
	return fmeta,nil
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
