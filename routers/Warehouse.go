package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services/warehouse"
	"../util"
	"strconv"
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
	val, err := strconv.ParseInt(key, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	reply, err := warehouse.GetOne(val)
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
	key := context.PostForm("fid")
	val, err := strconv.ParseInt(key, 10, 64)
	w.FId = val
	err = w.Check("Name", "Fid", "Group", "Version") // 参数检测
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

/**
	修改包信息
 */
func WarehouseModify(ctx *gin.Context) {
 	id := ctx.PostForm("id")
	val, err := strconv.ParseInt(id, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	old, err := warehouse.GetOne(val)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	w := warehouse.Default()
	w.Name = ctx.PostForm("name")
	w.Group = ctx.PostForm("group")
	w.Remarks = ctx.PostForm("remarks")
	w.Version = ctx.PostForm("version")
	reply, err := warehouse.Update(w, old)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, util.Success(reply))
}

 /**
 	删除文件包
  */
 func WarehouseDel(ctx *gin.Context) {
	 key := ctx.Param("key")
	 val, err := strconv.ParseInt(key, 10, 64)
	 if nil != err {
		 ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		 return
	 }
	 reply, err := warehouse.GetOne(val)
	 if nil != err {
		 ctx.JSON(http.StatusInternalServerError, util.Error(err))
		 return
	 }
	 reply.Status = false
	 _, err = reply.Modify()
	 if nil != err {
		 ctx.JSON(http.StatusInternalServerError, util.Error(err))
		 return
	 }
	 ctx.JSON(http.StatusOK, util.Success(reply))
 }