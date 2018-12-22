package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services/warehouse"
	"../util"
	"strconv"
)

func WarehouseList(context *gin.Context) {
	begin, limit := pageCondition(context)
	reply  := warehouse.GetList(strconv.Itoa(begin), strconv.Itoa(limit))
	context.JSON(http.StatusOK, util.Success(reply))
}
