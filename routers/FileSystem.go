package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../util"
	"../integrate/couchdb"
	"fmt"
	"time"
)

var (
	root, timeFormat string
	host, dbName string
)

func init() {
	root = "/tmp/filesystem"
	timeFormat = "2006-01-02 - 15:04:05"
	host = "mine:5984"
	dbName = "file-system"
}

type fsFile struct {
	name, savePath, contentType, key string
	uploadTime, size int64
	status bool
}

func (this *fsFile) generatorSavePath() string {
	return fmt.Sprintf("%s/%s", this.savePath, this.key)
}

func (this *fsFile) String() string {
	return fmt.Sprintf("name: %s savePaht: %s contentType: %s key: %s, uploadTime: %s, size: %d, status %v",
		this.name, this.savePath, this.contentType, this.key, time.Unix(this.uploadTime, 0).Format(timeFormat),
		this.size, this.status)
}

func (this *fsFile) generatorMap() map[string]interface{} {
	return map[string]interface{}{
		"name": this.name,
		"savePath": this.savePath,
		"contentType": this.contentType,
		"key": this.key,
		"uploadTime": this.uploadTime,
		"size": this.size,
		"status": this.status,
	}
}

func saveToCouchDB(object map[string]interface{}) (map[string]interface{}, error){
	return couchdb.Create(dbName, object)
}

/**
	文件上传
 */
func FsUpload(context *gin.Context) {
	file, err := context.FormFile("file")
	if nil != err {
		context.JSON(http.StatusBadRequest, util.Fail(400, "file not found."))
		return
	}
	fs := &fsFile{
		name: file.Filename,
		savePath: root,
		contentType: file.Header.Get("Content-Type"),
		key: util.GeneratorUUID(),
		uploadTime: time.Now().Unix(),
		size: file.Size,
		status: true,
	}
	object, err := saveToCouchDB(fs.generatorMap())
	if nil != err {
		context.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	err = context.SaveUploadedFile(file, fs.generatorSavePath())
	if nil != err {
		context.JSON(http.StatusInternalServerError, util.Fail(500, "io exception."))
		return
	}
	context.JSON(http.StatusOK, util.Success(object))
}