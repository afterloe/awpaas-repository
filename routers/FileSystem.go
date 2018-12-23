package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../util"
	"../integrate/soaClient"
	"fmt"
	"time"
	"strings"
)

var (
	root, timeFormat string
)

func init() {
	root = "F:/"
	timeFormat = "2006-01-02 - 15:04:05"
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

func saveToCouchDB(object map[string]interface{}) {
	reply, _ := soaClient.Call("GET", "192.168.3.21:5984", "/_uuids?count=1", nil, nil)
	id := reply["uuids"].([]interface{})[0]
	object["_id"] = id
	reply, _ = soaClient.Call("PUT", "192.168.3.21:5984", "/file-system/" + id.(string),
		strings.NewReader(soaClient.ObjectToJson(object)), map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		})
	fmt.Println(reply)
	if "not_found" == reply["error"]{
		soaClient.Call("PUT", "192.168.3.21:5984", "/file-system",
			nil, nil )
	}
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
	saveToCouchDB(fs.generatorMap())
	//context.SaveUploadedFile(file, fs.generatorSavePath())
	context.JSON(http.StatusOK, util.Success(fs.generatorMap()))
}