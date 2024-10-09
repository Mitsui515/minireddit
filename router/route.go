package route

import (
	"minireddit/controller"
	_ "minireddit/docs"
	"minireddit/logger"
	"minireddit/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func SetUpRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置为发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler)) // 添加swagger路由

	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler)              // 注册
	v1.POST("/login", controller.LoginHandler)                // 登录
	v1.POST("/refresh_token", controller.RefreshTokenHandler) // 刷新token

	v1.Use(middlewares.JWTAuthMiddleware()) // 使用JWT认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)           // 获取社区列表
		v1.GET("/community/:id", controller.CommunityDetailHandler) // 获取社区详情

		v1.POST("/post", controller.CreatePostHandler)       // 创建帖子
		v1.GET("/post/:id", controller.GetPostDetailHandler) // 获取帖子详情
		v1.GET("/posts", controller.GetPostListHandler)      // 获取帖子列表
		v1.GET("/posts2", controller.GetPostListHandler2)    // 根据时间或分数获取帖子列表

		v1.POST("/vote", controller.PostVoteHandler) // 投票

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
