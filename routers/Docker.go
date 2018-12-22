package routers

import (
	"github.com/gin-gonic/gin"
	"../services/warehouse"
	"../util"
	"net/http"
)

func DockerList(context *gin.Context) {
	reply  := warehouse.GetList()
	context.JSON(http.StatusOK, util.Success(reply))
}
