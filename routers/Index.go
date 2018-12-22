package routers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

/**
	路由列表
 */
func Execute(route *gin.RouterGroup) {
	route.GET("/remote/repository", RemoteList)
	route.GET("/warehouse/list", WarehouseList)
}

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
