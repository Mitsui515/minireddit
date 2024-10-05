package route

import (
	"minireddit/controller"
	"minireddit/logger"
	"minireddit/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置为发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler)              // 注册
	v1.POST("/login", controller.LoginHandler)                // 登录
	v1.POST("/refresh_token", controller.RefreshTokenHandler) // 刷新token

	v1.Use(middlewares.JWTAuthMiddleware()) // 使用JWT认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)           // 获取社区列表
		v1.GET("/community/:id", controller.CommunityDetailHandler) // 获取社区详情

		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
