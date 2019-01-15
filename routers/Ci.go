package routers

import (
	"github.com/gin-gonic/gin"
	"../util"
	"../services/warehouse"
	"net/http"
)

/**
	构建 - 获取ci类型
 */
func CmdList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, util.Success(warehouse.GetRegistryType()))
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
	fileType := ctx.PostForm("type")
	content := ctx.PostFormArray("content")
	cmd, err := warehouse.DefaultCmd(fileType, content...)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Error(err))
		return
	}
	old, err := warehouse.GetOne(key)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	reply, err := warehouse.AppendCI(old, cmd)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}