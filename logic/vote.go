package logic

import (
	"minireddit/dao/redis"
	"minireddit/models"
	"strconv"

	"go.uber.org/zap"
)

// 本项目使用简化版的投票分数
// 投一票加 432，86400/200 = 432 200票可以把帖子拉到一天的热度

// VoteForPost 为帖子投票
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
