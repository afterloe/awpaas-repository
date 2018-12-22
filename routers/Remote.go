package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services/remote"
	"../util"
)

func RemoteList(context *gin.Context) {
	reply  := remote.GetList()
	context.JSON(http.StatusOK, util.Success(reply))
}