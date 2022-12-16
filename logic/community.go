package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/model"

	"go.uber.org/zap"
)

// GetCommunityList 获取所有社区
func GetCommunityList() (communityList []*model.Community, err error) {
	// 查询数据库
	communityList, err = mysql.QueryCommunityList()
	if err != nil {
		zap.L().Error("mysql QueryCommunity() failed", zap.Error(err))
		return
	}
	return
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(id int) (community *model.CommunityDetail, err error) {
	community, err = mysql.GetCommunity(id)
	if err != nil {
		zap.L().Error("mysql.GetCommunity() failed", zap.Error(err))
		return
	}
	return
}
