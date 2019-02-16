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

	// 上传镜像
	route.PUT("/file", FsUpload)
	route.GET("/download/:key", FsDownload) // 下载

	route.GET("/remote/repository", RemoteList) // 远程 - 镜像列表

	route.POST("/warehouse/file/:key", WarehouseLoad) // 文件上传
	route.GET("/remote/repository/:name/version") // 远程 - 指定镜像版本列表
	route.GET("/remote/detail/:name/:version") // 远程 - 镜像详情

	// 镜像管理
	route.GET("/warehouse", WarehouseList) // 镜像管理 - 镜像查询
	route.GET("/warehouse/:key", WarehouseOne) // 镜像管理 - 查看详情
	route.PUT("/warehouse", WarehouseAppend) // 镜像管理 - 镜像创建
	route.POST("/warehouse", WarehouseModify) // 镜像管理 - 镜像信息修改
	route.DELETE("/warehouse/:key", WarehouseDel) // 镜像管理 - 镜像删除
	route.GET("/export/:key") // 镜像管理 - 导出

	// 构建
	route.GET("/ci/type", CmdList) // 构建类型列表
	route.GET("/ci/warehouse/:key", CmdGet) // 构建 - 查询构建信息
	route.PUT("/ci/warehouse/:key", CmdBuilder) // 构建 - 创建构建命令
	route.DELETE("/ci/warehouse/:key") // 构建 - 删除构建命令
	route.POST("/ci/warehouse/:key", CmdCi) // 构建 - 执行构建命令

	// 镜像运行
	route.PUT("/run/:key") // 执行镜像
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
