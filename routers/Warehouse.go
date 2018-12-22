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
	reply  := warehouse.GetList(strconv.Itoa(begin), strconv.Itoa(limit))
	context.JSON(http.StatusOK, util.Success(reply))
}

/**
	添加包
 */
func WarehouseAppend(context *gin.Context) {
	//name := context.PostForm("name")
	//group := context.PostForm("group")
	//remarks := context.PostForm("remarks")
	//version := context.PostForm("version")

}
