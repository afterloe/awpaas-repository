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
	begin, limit := pageCondition(context)
	reply := warehouse.GetList(strconv.Itoa(begin), strconv.Itoa(limit))
	context.JSON(http.StatusOK, util.Success(reply))
}

/**
	添加包
 */
func WarehouseAppend(context *gin.Context) {
	module := &warehouse.Module {
		Name: context.PostForm("name"),
		Group: context.PostForm("group"),
		Remarks: context.PostForm("remarks"),
		Version : context.PostForm("version"),
		Fid: context.PostForm("fid"),
	}
	err := module.Check("Name", "Fid", "Group", "Version") // 参数检测
	if nil != err {
		context.JSON(http.StatusBadRequest, util.Error(err))
		return
	}
	reply, err := module.AppendToRemote()
	if nil != err {
		context.JSON(http.StatusBadRequest, util.Error(err))
		return
	}
	context.JSON(http.StatusOK, util.Success(reply))
}
