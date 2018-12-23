package routers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"../util"
	"../config"
)

/**
	路由列表
 */
func Execute(route *gin.RouterGroup) {
	route.POST("/fs/upload", FsUpload)
	route.GET("/remote/repository", RemoteList)
	route.GET("/warehouse", WarehouseList)
	route.PUT("/warehouse", WarehouseAppend)
}

/**
	描述信息
 */
func Info(context *gin.Context) {
	info := config.Get("info").(map[string]interface{})
	context.JSON(http.StatusOK, util.Success(info))
}

/**
	分页组件
 */
func pageCondition(context *gin.Context) (int, int) {
	begin, err := strconv.Atoi(context.DefaultQuery("bg", "0"))
	if nil != err {
		begin = 0
	}
	end, err := strconv.Atoi(context.DefaultQuery("ed", "10"))
	if nil != err {
		end = 10
	}
	limit := end - begin
	if 0 >= limit {
		limit = 10
	}
	return begin, limit
}
