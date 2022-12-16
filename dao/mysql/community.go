package mysql

import (
	"awesomeProject/model"
)

func QueryCommunityList() (communityList []*model.Community, err error) {
	// sql语句
	sql := `select community_id,community_name from community where id > ?`
	err = db.Select(&communityList, sql, 0)
	if err != nil {
		return nil, err
	}
	return
}

func GetCommunity(id int) (community *model.CommunityDetail, err error) {
	sql := `select community_id, community_name, introduction from community where community_id = ?`
	community = new(model.CommunityDetail)
	err = db.Get(community, sql, id)
	if err != nil {
		return
	}
	return
}
