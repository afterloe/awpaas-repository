package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services/warehouse"
	"../util"
)

/**
	查询创建的包
 */
func WarehouseList(context *gin.Context) {
	reply := warehouse.GetList(pageCondition(context))
	context.JSON(http.StatusOK, util.Success(reply))
}

/**
	查询包详情
*/
func WarehouseOne(ctx *gin.Context) {
	key := ctx.Param("key")
	if 32 > len(key) {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	reply, err := warehouse.GetOne(key)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(*reply))
}

/**
	添加包
 */
func WarehouseAppend(context *gin.Context) {
	w := warehouse.Default()
	w.Name = context.PostForm("name")
	w.Group = context.PostForm("group")
	w.Remarks = context.PostForm("remarks")
	w.Version = context.PostForm("version")
	w.Fid = context.PostForm("fid")
	err := w.Check("Name", "Fid", "Group", "Version") // 参数检测
	if nil != err {
		context.JSON(http.StatusBadRequest, util.Error(err))
		return
	}
	object, err := w.SaveToDB()
	if nil != err {
		context.JSON(http.StatusBadRequest, util.Error(err))
		return
	}
	context.JSON(http.StatusOK, util.Success(object))
}
