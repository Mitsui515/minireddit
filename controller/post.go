package controller

import (
	"minireddit/logic"
	"minireddit/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子的处理函数
// @Summary 创建帖子
// @Description 创建一个新的帖子
// @Tags 帖子
// @Accept  json
// @Produce  json
// @Param   post  body  models.Post  true  "帖子内容"
// @Success 200 {object} _ResponsePostList
// @Failure 400 {object} _ResponsePostList
// @Failure 500 {object} _ResponsePostList
// @Router /posts [post]
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON failed", zap.Any("err", err))
		zap.L().Error("CreatePost with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从请求中获取发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 处理获取帖子详情的请求
// @Summary 获取帖子详情
// @Description 根据帖子ID获取帖子详情
// @Tags 帖子相关接口
// @Accept  json
// @Produce  json
// @Param   id     path    int64     true        "帖子ID"
// @Success 200 {object} _ResponsePostList "成功获取帖子详情"
// @Failure 400 {object} _ResponsePostList "请求参数错误"
// @Failure 500 {object} _ResponsePostList "服务器内部错误"ResponsePostList
// @Router /posts/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数及参数校验 (从URL中获取帖子ID)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据帖子ID获取帖子详情
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 处理获取帖子列表的请求。
// 它从请求上下文中提取分页参数，从逻辑层获取帖子列表，
// 并将适当的响应发送回客户端。
//
// @Summary 获取帖子列表
// @Description 获取分页的帖子列表
// @Tags 帖子
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param size query int true "每页大小"
// @Success 200 {object} _ResponsePostList "帖子列表"
// @Failure 500 {object} _ResponsePostList "服务器内部错误"
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 从请求中获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数(query string): /api/v1/posts2?page=1&size=10&order=time
	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := logic.GetPostListNew(p) // 更新后的逻辑
	if err != nil {
		zap.L().Error("logic.GetPostListNew failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据社区去查询帖子列表
// func GetCommunityPostListHandler(c *gin.Context) {
// 	p := &models.ParamCommunityPostList{
// 		ParamPostList: &models.ParamPostList{
// 			Page:  1,
// 			Size:  10,
// 			Order: models.OrderTime,
// 		},
// 	}
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
// 		ResponseError(c, CodeInvalidParam)
// 		return
// 	}
// 	// 获取数据
// 	data, err := logic.GetCommunityPostList(p)
// 	if err != nil {
// 		zap.L().Error("logic.GetCommunityPostList failed", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	// 返回响应
// 	ResponseSuccess(c, data)
// }
