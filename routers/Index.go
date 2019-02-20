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

	// 文件管理
	route.PUT("/file", FsUpload) // 文件上传
	route.DELETE("/file") // 文件删除
	route.GET("/file") // 文件详情
	route.GET("/file/list") // 已上传文件列表
	route.GET("/file/download/:key", FsDownload) // 文件下载

	// 远程管理
	route.GET("/remote/repository", RemoteList) // 获取远程镜像列表
	route.GET("/remote/repository/:name/version") // 获取远程指定镜像版本列表
	route.GET("/remote/detail/:name/:version") // 获取远程指定镜像详情

	// 镜像管理
	route.PUT("/warehouse", WarehouseAppend) //镜像创建
	route.GET("/warehouse", WarehouseList) // 镜像查询
	route.POST("/warehouse", WarehouseModify) // 镜像管理 - 镜像信息修改
	route.GET("/warehouse/:key", WarehouseOne) // 镜像管理 - 查看详情
	route.DELETE("/warehouse/:key", WarehouseDel) // 镜像管理 - 镜像删除
	route.GET("/export/:key") // 镜像管理 - 导出

	// 常量字典
	route.GET("/dictionaries/ci", CmdList) // 构建类型列表

	// 构建
	route.GET("/ci/warehouse/:key", CIList) // 查询构建信息
	route.PUT("/ci/warehouse/:key", CmdBuilder) // 创建构建命令
	route.GET("/ci/item/:key", CmdGet) // 获取构建命令详情
	route.POST("/ci/item/:key", RunCI) // 执行构建命令
	route.GET("/ci/history", CIHistory) // 执行历史
	route.GET("/ci/history/:key", CIHistoryDetail) // 详细执行历史
	route.GET("/ci/detail/:key", ) // 查询ci执行结果
	route.DELETE("/ci/warehouse/:key") // 删除构建命令

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
