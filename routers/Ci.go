package routers

import (
	"github.com/gin-gonic/gin"
	"../util"
	"net/http"
)

func CmdList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, util.Success([...]string{"code", "image", "tar"}))
}

/**
	构建 - 查询构建信息
 */
func CmdGet(ctx *gin.Context) {
	key := ctx.Param("key")
	if 32 > len(key) {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
}

/**
	构建 - 创建构建命令
 */
func CmdBuilder(ctx *gin.Context) {
	key := ctx.Param("key")
	if 32 > len(key) {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}

}