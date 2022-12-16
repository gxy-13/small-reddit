package routers

import (
	"awesomeProject/controller"
	"awesomeProject/logger"
	"awesomeProject/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	//gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/ping", middleware.JwtAuth(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r1 := r.Group("/api/v1")

	// 注册
	r1.POST("/signup", controller.SignUp)
	// 登陆
	r1.POST("/login", controller.Login)

	// 使用JWT中间件校验权限
	r1.Use(middleware.JwtAuth())

	{
		// 获取社区分类
		r1.GET("/community", controller.GetCommunity)
		// 获取社区详情
		r1.GET("/community/:id", controller.GetCommunityDetail)
		// 发布帖子
		r1.POST("/post", controller.CreatePost)
		// 获取帖子详情
		r1.GET("/post/:id", controller.GetPostDetail)
		//	根据时间或者分数获取帖子列表
		r1.GET("/posts2", controller.GetPostList)
		// 投票
		r1.POST("/vote", controller.PostVote)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": 404,
		})
	})
	return r
}
