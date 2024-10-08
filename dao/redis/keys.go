package redis

// Redis Key
// Redis Key注意使用命名空间，方便查询和拆分，避免key冲突
const (
	KeyPrefix              = "minireddit:"
	KeyPostTimeZSet        = "post:time"   // ZSet; 帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  // ZSet; 帖子及投票得分
	KeyPostVotedZSetPrefix = "post:voted:" // ZSet; 记录用户及投票类型; 参数是postID
)

// getRedisKey 获取带有统一前缀的key
func getRedisKey(key string) string {
	return KeyPrefix + key
}
