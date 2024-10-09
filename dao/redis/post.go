package redis

import (
	"minireddit/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 1. 确定查询的索引范围
	start := (page - 1) * size
	end := start + size - 1
	// 2. ZREVRANGE 按分数从大到小的数据查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

// GetPostIDsInOrder 根据用户请求参数获取要查询的post id列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis中获取id
	// 1. 根据用户请求中携带的order参数确定要查询的key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投票数
func GetPostVoteData(ids []string) (data []int64, err error) {
	// 使用pipeline一次发送多个命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 根据社区ID获取社区下的帖子id列表
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用ZINTERSTORE把分区的帖子set与帖子分数的zset做交集，生成一个新的zset
	// 针对新的zset进行分页查询

	// 社区的key
	cKey := KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID))

	// 利用缓存key减少zinterstore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.TxPipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, orderKey, cKey)
		// 设置过期时间
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在，直接查询
	return getIDsFromKey(key, p.Page, p.Size)
}
