package routes

import (
	"minireddit/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	return r
}
