package routers

import (
	"github.com/gin-gonic/gin"
	"../util"
	"../services/ciTool"
	"net/http"
	"strconv"
)

/**
	构建 - 获取ci类型
 */
func CmdList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, util.Success(ciTool.GetRegistryType()))
}

/**
	查询构建详情
 */
func CmdGet(ctx *gin.Context) {
	key := ctx.Param("key")
	val, err := strconv.ParseInt(key, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	reply, err := ciTool.GetOne(val)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}

/**
	构建 - 查询构建信息列表
 */
func CIList(ctx *gin.Context) {
	key := ctx.Param("key")
	val, err := strconv.ParseInt(key, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	reply, err := ciTool.CIList(val)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}

/**
	构建 - 添加构建命令
 */
func CmdBuilder(ctx *gin.Context) {
	key := ctx.Param("key")
	val, err := strconv.ParseInt(key, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	fileType := ctx.PostForm("type")
	content := ctx.PostFormArray("content")
	cmd, err := ciTool.DefaultCmd(val, fileType, content...)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Error(err))
		return
	}
	reply, err := ciTool.AppendCI(val, fileType, cmd)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}

/**
	构建 - 执行构建命令
 */
func RunCI(ctx *gin.Context) {
	key := ctx.Param("key")
	val, err := strconv.ParseInt(key, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	reply, err := ciTool.Run(val)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}

func CIHistory(ctx *gin.Context) {
	reply, err := ciTool.History(pageCondition(ctx))
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}

func CIDetail(ctx *gin.Context) {
	key := ctx.Param("key")
	if 32 != len(key) {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	reply, err := ciTool.GetDetail(key)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.Status(200)
	ctx.Header("Content-Type", "text")
	ctx.Header("Content-Disposition", "attachment;filename=" + string([]byte(reply["name"].(string))))
	ctx.Header("Content-Length", string(reply["size"].(int64)))
	ctx.File(reply["path"].(string))
}

func CIHistoryDetail(ctx *gin.Context) {
	key := ctx.Param("key")
	val, err := strconv.ParseInt(key, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	reply, err := ciTool.CIHistoryDetail(val)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}
