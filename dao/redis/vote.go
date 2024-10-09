package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

/* 投票的几种情况
direction=1时，有两种情况:
	1. 之前没有投过票，现在投赞成票 --> 更新分数和投票纪录 差值的绝对值: 1 +432
	2. 之前投反对票，现在改投赞成票 --> 更新分数和投票纪录 差值的绝对值: 2 +864
direction=0时，有两种情况:
	1. 之前投赞成票，现在要取消投票 --> 更新分数和投票纪录 差值的绝对值: 1 -432
	2. 之前投反对票，现在要取消投票 --> 更新分数和投票纪录 差值的绝对值: 1 +432
direction=-1时，有两种情况:
	1. 之前没有投过票，现在投反对票 --> 更新分数和投票纪录 差值的绝对值: 1 -432
	2. 之前投赞成票，现在改投反对票 --> 更新分数和投票纪录 差值的绝对值: 2 -864

投票的限制:
每个帖子自发表之日起一个星期内允许投票，超过一个星期不允许投票
	1. 到期之后将redis中的赞成票数和反对票数存到mysql中
	2. 到期之后删除KeyPostVotedZSetPrefix
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 把帖子ID加到社区的set里面
	pipeline.SAdd(getRedisKey(KeyCommunitySetPrefix+strconv.Itoa(int(communityID))), postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2和3需要放在一个事务中
	// 2. 更新分数
	// 先查看用户给该帖子的投票记录
	oldValue := client.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	// 如果这次的投票和上次的投票一样，就提示不允许重复投票
	if oldValue == value {
		return ErrVoteRepeated
	}
	var op float64
	if value > oldValue {
		op = 1
	} else if value < oldValue {
		op = -1
	}
	diff := math.Abs(oldValue - value) // 计算分数差
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID).Result()
	// 3. 记录用户投票数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value,  // 赞成票还是反对票
			Member: userID, // 用户ID
		})
	}
	_, err := pipeline.Exec()
	return err
}
