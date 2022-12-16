package redis

// 固定的不变的字符串定义为常量
// 定义key的时候使用命名空间形式，方便查找和拆分
const (
	KeyPrefix             = "Free:"       // 前缀
	KeyPostTimeZSet       = "post:time"   //帖子及时间
	KeyPostScoreZSet      = "post:score"  //帖子及分数
	KeyPostVotedPrefix    = "post:voted:" //不同帖子不同用户的投票状态， 比如post:voted:11343  表示用户11343 帖子的投标状态
	KeyCommunitySetPrefix = "community:"  //set 保存每个分区下帖子的id
)

// 给redis key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
