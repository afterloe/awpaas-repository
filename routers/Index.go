package routers

import (
	"github.com/gin-gonic/gin"
)

/**
	路由列表
 */
func Execute(route *gin.RouterGroup) {
	route.GET("/docker/repository", DockerList)
}