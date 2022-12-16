package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

// scorePerVote 表示每一个票加减500分  duration 表示一个帖子可以投票的时长为7天
const (
	scorePerVote = 500
	duration     = 86400 * 7
)

var (
	ErrPostTimeExpired = errors.New("超出投票时间")
	ErrVoteRepeated    = errors.New("不允许重复投票")
)

/*
	direction 有三种情况1，0 -1。每种情况又有两种子情况
	1:
		之前没投票，现在投赞成    0 -> 1  500
		之前投反对，现在投赞成    -1 -> 1 500 * 2
	0:
		之前投反对, 现在取消投票  -1 -> 0  500
		之前是赞成，现在取消投票  1 -> 0  500
	-1:
		之前是赞成，现在投反对   1 -> -1  -500 * 2
		之前没投票，现在投反对   0 -> -1 - 500
*/
// PostVote
func PostVote(pid, uid string, direction float64) (err error) {
	//// 1.判断投票限制
	//// 判断帖子是否还在有效期内，获取帖子的发布时间
	//postTime := rdb.ZScore(KeyPrefix+KeyPostTimeZSet, v.PostID).Val()
	//if float64(time.Now().Unix())-postTime > duration {
	//	zap.L().Error("post is expired", zap.Error(err))
	//	return ErrPostTimeExpired
	//}
	//// 第2步和第3步也需要放到一个pipeline事务中去
	//// 2.计算最新的分数
	//// 通过上面分析投票状态可以得出，当新direction > 旧direction，分数减小 反之是增加
	//// 比较direction
	//var dir float64
	//newRes := rdb.ZScore(KeyPostVotedPrefix+v.PostID, uid).Val()
	//// 如果这一次投票的值和之前保存的值一致，就提示不允许投票
	//if dir == newRes {
	//	return ErrVoteRepeated
	//}
	//nDir, err := strconv.Atoi(v.Direction)
	//if err != nil {
	//	return
	//}
	//if float64(nDir) > newRes {
	//	dir = 1
	//} else {
	//	dir = -1
	//}
	//// 计算新投票和老投票之间的差值
	//diff := math.Abs(float64(nDir) - newRes)
	//// 计算新分数
	//pipeline := rdb.TxPipeline()
	//pipeline.ZIncrBy(KeyPostScoreZSet, dir*diff*scorePerVote, v.PostID)
	//if ErrPostTimeExpired != nil {
	//	return err
	//}
	//// 3.将帖子用户投票的状态更新
	//// 说明用户是取消了投票
	//if dir == 0 {
	//	pipeline.ZRem(KeyPostVotedPrefix+v.PostID, uid)
	//} else {
	//	pipeline.ZAdd(KeyPostVotedPrefix+v.PostID, redis.Z{
	//		Score:  dir, // 用户新投的是赞成票还是反对票
	//		Member: uid,
	//	})
	//}
	//_, err = pipeline.Exec()
	//return
	// 1. 判断投票限制 从redis取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), pid).Val()
	if float64(time.Now().Unix())-postTime > duration {
		return ErrPostTimeExpired
	}
	// 2 和 3 需要放到一个pipeline事务中操作

	// 2.更新帖子的分数
	// 先查看当前用户给当前帖子的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedPrefix+pid), uid).Val()

	// 更新 如果这一次投票的值和之前保存的值一直，就提示不允许重复投票
	if direction == ov {
		return ErrVoteRepeated
	}
	var op float64
	if direction > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - direction) // 计算两次投票的差值
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, pid)

	// 3 记录用户为该帖子投票的数据
	if direction == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedPrefix+pid), uid)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedPrefix+pid), redis.Z{
			Score:  direction, // 赞成还是反对
			Member: uid,
		})
	}
	_, err = pipeline.Exec()
	return err
}
