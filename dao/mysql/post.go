package mysql

import (
	"awesomeProject/model"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CreatePost 发布帖子
func CreatePost(p *model.PostDetail) (err error) {
	fmt.Printf("%#v\n", p)
	sql := `insert into post(post_id,title,content,author_id,community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sql, p.PostID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.CreatePost() failed", zap.Error(err))
		return
	}
	return
}

// GetPostDetail 获取帖子详情
func GetPostDetail(pid int64) (post *model.PostDetail, err error) {
	post = new(model.PostDetail)
	fmt.Printf("%d\n", pid)
	sql := `select post_id,content,title,community_id,author_id, status,create_time,update_time from post where post_id = ?`
	err = db.Get(post, sql, pid)
	if err != nil {
		return
	}
	return
}

// GetAllPosts 获取帖子列表
func GetAllPosts(page, size int) (postList []*model.PostDetail, err error) {
	sql := `select 
    		post_id,title,content,author_id,community_id,create_time,update_time 
			from post
			ORDER BY create_time DESC 
			limit ?,?		
	`
	err = db.Select(&postList, sql, page-1, size)
	if err != nil {
		return
	}
	return
}

// GetPostListByIDs 根据id列表查询帖子
func GetPostListByIDs(ids []string) (postList []*model.PostDetail, err error) {
	sql := `select post_id, title, content, author_id, community_id, create_time from post
			where post_id in (?)
			order by FIND_IN_SET(post_id,?)
			`
	query, args, err := sqlx.In(sql, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	query = db.Rebind(query)

	err = db.Select(&postList, query, args...)
	return
}
