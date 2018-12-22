package routers

import (
	"github.com/gin-gonic/gin"
	"../integrate/soaClient"
	"../util"
	"net/http"
)

func DockerList(context *gin.Context) {
	reply, err := soaClient.Call("GET",  "192.168.1.3", "/v2/_catalog", nil, nil)
	if nil != err {
		context.JSON(http.StatusInternalServerError, util.Fail(500, "call fall"))
		return
	}
	context.JSON(http.StatusOK, util.Success(reply))
}
