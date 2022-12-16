package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/dao/redis"
	"awesomeProject/model"
	"awesomeProject/utils/snowflake"
	"fmt"

	"go.uber.org/zap"
)

// CreatePost 发帖
func CreatePost(p *model.PostDetail) (err error) {
	// 生成PID
	p.PostID = snowflake.GenID()
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}
	err = redis.CreatePost(p.PostID, p.CommunityID)

	return
}

// GetPostDetail 获取帖子详情
func GetPostDetail(pid int64) (post *model.ApiPostDetail, err error) {
	p := new(model.PostDetail)
	p, err = mysql.GetPostDetail(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostDetail(pid)", zap.Error(err))
		return
	}
	// 获取用户信息
	username, err := mysql.GetUserByID(p.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(p.AuthorID)", zap.Error(err))
		return
	}
	community, err := mysql.GetCommunity(p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunity(p.CommunityID) failed", zap.Error(err))
		return
	}
	post = &model.ApiPostDetail{
		Username:        username,
		PostDetail:      p,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int) (postList []*model.ApiPostDetail, err error) {
	// 需要获取详细数据，所以需要对列表中的每一个post 都进行一次查询详细信息
	posts, err := mysql.GetAllPosts(page, size)
	if err != nil {
		zap.L().Error("mysql.GetAllPosts(page, size) failed", zap.Error(err))
		return
	}
	// 初始化所有帖子详情列表，长度就是帖子个数
	postList = make([]*model.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 获取用户信息
		username, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			continue
		}
		community, err := mysql.GetCommunity(post.CommunityID)
		if err != nil {
			continue
		}
		// 将每一个帖子详情信息都追加进postList
		postList = append(postList, &model.ApiPostDetail{
			Username:        username,
			PostDetail:      post,
			CommunityDetail: community,
		})
	}
	return
}

// GetPostListPlus 根据参数动态获取帖子列表
func GetPostListPlus(p *model.ParamPostList) (postList []*model.ApiPostDetail, err error) {
	// 去redis 查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder success, but no data")
		return
	}
	// 根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	fmt.Printf("%#v\n", voteData)
	if err != nil {
		return
	}
	// 将帖子的作者和社区信息查询并填充
	for idx, post := range posts {
		// 获取用户信息
		username, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			continue
		}
		community, err := mysql.GetCommunity(post.CommunityID)
		if err != nil {
			continue
		}
		// 将每一个帖子详情信息都追加进postList
		postList = append(postList, &model.ApiPostDetail{
			Username:        username,
			VoteNum:         voteData[idx],
			PostDetail:      post,
			CommunityDetail: community,
		})
		fmt.Println(voteData[idx])
	}
	return
}

// PostVote 投票逻辑
func PostVote(v *model.VoteParam, userID int64) (err error) {
	return redis.PostVote(v.PostID, fmt.Sprint(userID), float64(v.Direction))
}

func GetCommunityPostList(p *model.ParamCommunityPostList) (postList []*model.ApiPostDetail, err error) {
	// 去redis 查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder success, but no data")
		return
	}
	// 根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	// 将帖子的作者和社区信息查询并填充
	for idx, post := range posts {
		// 获取用户信息
		username, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			continue
		}
		community, err := mysql.GetCommunity(post.CommunityID)
		if err != nil {
			continue
		}
		// 将每一个帖子详情信息都追加进postList
		postList = append(postList, &model.ApiPostDetail{
			Username:        username,
			VoteNum:         voteData[idx],
			PostDetail:      post,
			CommunityDetail: community,
		})
	}
	return
}
