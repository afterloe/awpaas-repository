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
