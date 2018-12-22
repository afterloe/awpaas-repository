package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../util"
)

/**
	文件上传
 */
func FsUpload(context *gin.Context) {
	file, err := context.FormFile("file")
	if nil != err {
		context.JSON(http.StatusBadRequest, util.Fail(400, "file not found."))
		return
	}
	context.SaveUploadedFile(file, "C:/Users/afterloe/Desktop/" + util.GeneratorUUID())
	context.JSON(http.StatusOK, util.Success(""))
}