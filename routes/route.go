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

	r.POST("/signup", controller.SignUpHandler)              // 注册
	r.POST("/login", controller.LoginHandler)                // 登录
	r.POST("/refresh_token", controller.RefreshTokenHandler) // 刷新token

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
