package model

type Community struct {
	CommunityID   int    `db:"community_id" json:"id"`
	CommunityName string `db:"community_name" json:"name"`
}

type CommunityDetail struct {
	CommunityID   int    `db:"community_id" json:"community_id"`
	CommunityName string `db:"community_name" json:"community_name"`
	Introduction  string `db:"introduction" json:"introduction"`
}

type ParamCommunityPostList struct {
	*ParamPostList
	CommunityID int `json:"community_id" form:"community_id"`
}
