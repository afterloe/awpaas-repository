package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services/warehouse"
	"../util"
	"strconv"
)

func WarehouseList(context *gin.Context) {
	var (
		begin_str = context.DefaultQuery("bg", "0")
		end_str = context.DefaultQuery("ed", "10")
	)
	begin, err := strconv.Atoi(begin_str)
	if nil != err {
		begin = 0
	}
	end, err := strconv.Atoi(end_str)
	if nil != err {
		end = 10
	}
	limit := end - begin
	if 0 >= limit {
		limit = 10
	}
	reply  := warehouse.GetList(strconv.Itoa(begin), strconv.Itoa(limit))
	context.JSON(http.StatusOK, util.Success(reply))
}
