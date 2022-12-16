package controller

import (
	"awesomeProject/logic"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// GetCommunity 获取社区分类
func GetCommunity(c *gin.Context) {
	// 获取参数和参数校验
	// 业务逻辑 需要获取community_id community_name 存放在切片中
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(c *gin.Context) {
	// 获取参数
	communityID := c.Param("id")
	cid, err := strconv.Atoi(communityID)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	community, err := logic.GetCommunityDetail(cid)
	if err != nil {
		zap.L().Error("logic.getCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, community)
}
