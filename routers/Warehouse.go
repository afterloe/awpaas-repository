package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services/warehouse"
	"../services/fileSystem"
	"../util"
)

/**
	查询创建的包
 */
func WarehouseList(context *gin.Context) {
	reply := warehouse.GetList(pageCondition(context))
	context.JSON(http.StatusOK, util.Success(reply))
}

func WarehouseLoad(ctx *gin.Context) {
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
	file, err := ctx.FormFile("file")
	if nil != err {
		ctx.SecureJSON(http.StatusBadRequest, util.Fail(400, "file not found."))
		return
	}
	fs := fileSystem.Default(file.Filename, file.Header.Get("Content-Type"), file.Size)
	reply.PackInfo = *fs
	if nil != err {
		ctx.SecureJSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	err = ctx.SaveUploadedFile(file, fs.GeneratorSavePath())
	reply.Modify()
	if nil != err {
		ctx.SecureJSON(http.StatusInternalServerError, util.Fail(500, "io exception."))
		return
	}
	ctx.SecureJSON(http.StatusOK, util.Success(true))
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

/**
	修改包信息
 */
func WarehouseModify(ctx *gin.Context) {
 	id := ctx.PostForm("id")
 	if "" == id {
		ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		return
	}
	old, err := warehouse.GetOne(id)
	if nil != err {
		ctx.JSON(http.StatusInternalServerError, util.Error(err))
		return
	}
	w := warehouse.Default()
	w.Name = ctx.PostForm("name")
	w.Group = ctx.PostForm("group")
	w.Remarks = ctx.PostForm("remarks")
	w.Version = ctx.PostForm("version")
	w.Fid = ctx.PostForm("fid")
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
	 if 32 > len(key) {
		 ctx.JSON(http.StatusBadRequest, util.Fail(400, "参数错误"))
		 return
	 }
	 reply, err := warehouse.GetOne(key)
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