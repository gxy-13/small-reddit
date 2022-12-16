package redis

import (
	"awesomeProject/model"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func CreatePost(postID int64, communityID int) error {
	pipeline := rdb.TxPipeline()
	// 帖子时间 在redis中存放帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数 在redis中创造帖子的同时创造一个和时间相同的分数，这样每次投票都会让帖子存活时间变长
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(communityID))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 确定查询索引的起始点
	start := (page - 1) * size
	end := start + size - 1
	// ZRevRange 查询 按分数从大到小查询指定数量分数
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *model.ParamPostList) ([]string, error) {
	// 从redis获取ID，根据请求中的order字段决定要查询的key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == model.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//数量大的时候很耗时，可以使用pipeline
	//data = make([]int64, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedPrefix + id)
	//	// 查找key中分数是1的元素的数量，统计每个帖子的赞成票数量
	//	v1 := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}
	fmt.Printf("%#v\n", ids)
	// 使用pipeline 一次发送多条命令，减少rtt
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + "" + id)
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

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *model.ParamCommunityPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	// 使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	// 针对新的zset按之前的逻辑取数据
	// 利用缓存key 减少 zinterstore执行的次数
	key := orderKey + strconv.Itoa(p.CommunityID)
	if rdb.Exists(orderKey).Val() < 1 {
		// 不存在需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, getRedisKey(KeyCommunitySetPrefix+strconv.Itoa(p.CommunityID)), orderKey) // 计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在就直接根据key查询ids

	return getIDsFromKey(key, p.Page, p.Size)
}
