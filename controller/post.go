package controller

import (
	"awesomeProject/logic"
	"awesomeProject/model"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CreatePost 发布帖子
func CreatePost(c *gin.Context) {
	// 获取参数，参数校验
	p := new(model.PostDetail)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("post with invalid param", zap.Error(err))
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Translate(trans))
	}
	zap.L().Info("post title", zap.String("post.title", p.Title))
	// 业务逻辑
	// 先获取当前userID
	uid, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("user no login")
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = uid
	err = logic.CreatePost(p)
	if err != nil {
		zap.L().Error("logic.CreatePost() failed")
		if err == sql.ErrNoRows {
			ResponseError(c, CodeServerBusy)
			return
		}
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetail 获取帖子详情
func GetPostDetail(c *gin.Context) {
	// 获取参数
	pid := c.Param("id")
	postID, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务逻辑
	post, err := logic.GetPostDetail(postID)
	if err != nil {
		zap.L().Error("logic.GetPostDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, post)
}

// GetPostList 根据前端传递的参数动态获取帖子列表，
// 按创建时间排序，或者按照分数排序
// 获取参数
// 去redis获取id列表
// 根据id查询帖子详情
func GetPostList(c *gin.Context) {
	// 获取参数， 参数校验
	//page, size := GetPageInfo(c)
	// 获取帖子使用的是get请求，get请求的参数是写在url上的并不是JSON格式，三个参数，page ，size，order 可以使用一个结构体
	// 利用结构体实现默认值
	p := &model.ParamPostList{
		Page:  1,
		Size:  10,
		Order: model.OrderTime, // 减少magic string
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostList invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务逻辑
	postList, err := logic.GetPostListPlus(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, postList)
}

// GetPageInfo 获取分页信息
func GetPageInfo(c *gin.Context) (page, size int) {
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		// 如果没取到就给一个默认值
		page = 1
	}
	sizeStr := c.Query("size")
	size, err = strconv.Atoi(sizeStr)
	if err != nil {
		// 没取到就给一个默认值
		size = 3
	}
	return
}

// PostVote 投票
func PostVote(c *gin.Context) {
	// 获取参数，校验参数，投票需要 user_id 可以通过上下文获取， post_id 必要， 投票的种类也是必要的
	p := new(model.VoteParam)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Vote with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}
	// 业务逻辑
	// 先获取userID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	logic.PostVote(p, userID)
	ResponseSuccess(c, CodeSuccess)
}

// GetCommunityPostList 根据社区查询帖子
func GetCommunityPostList(c *gin.Context) {
	p := &model.ParamCommunityPostList{
		ParamPostList: &model.ParamPostList{
			Page:  1,
			Size:  10,
			Order: model.OrderTime, // 减少magic string
		},
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostList invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务逻辑
	postList, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, postList)
}
