package model

import "time"

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// PostDetail 帖子信息
type PostDetail struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	CommunityID int       `json:"community_id" db:"community_id" binding:"required"`
	Status      int       `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
}

// ApiPostDetail 用于接口的post信息
type ApiPostDetail struct {
	Username string `json:"username"`
	VoteNum  int64  `json:"vote_num"`
	*PostDetail
	*CommunityDetail `json:"community"`
}

// ParamPostList 获取帖子列表的三个参数,get请求，tag使用的是form
type ParamPostList struct {
	Page  int64  `form:"page" json:"page"`
	Size  int64  `form:"size" json:"size"`
	Order string `form:"order" json:"order"`
}
