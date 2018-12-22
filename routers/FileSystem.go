package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../util"
)

/**
	描述信息
 */
func FsUpload(context *gin.Context) {
	context.JSON(http.StatusOK, util.Success(""))
}