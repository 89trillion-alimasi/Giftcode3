package router

import (
	"GiftCode2/controller"
	"github.com/gin-gonic/gin"
)

// InitRouter 负责初始化 web 框架的路由器
func InitRouter() *gin.Engine {
	// 创建路由器
	route := gin.Default()

	// 添加路由
	route.POST("/create_gift_code", controller.CreateGiftCode)
	route.GET("/query_gift_code", controller.QueryGiftCode)
	route.POST("/verify_gift_code", controller.VerifyGiftCode)

	return route
}
