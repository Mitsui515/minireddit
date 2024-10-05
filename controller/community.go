package controller

import (
	"minireddit/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 获取社区列表
func CommunityHandler(c *gin.Context) {
	// 1. 查询所有社区(community_id, community_name)，以列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易暴露系统内部错误
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 获取社区详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据id查询社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}
