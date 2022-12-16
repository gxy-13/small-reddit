package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseError 错误请求响应，返回没有数据，只有code 和 msg
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  code.getMsg(),
		"data": nil,
	})
}

// ResponseSuccess  成功请求响应， 返回code msg data
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  CodeSuccess.getMsg(),
		"data": data,
	})
}

// ResponseErrorWithMsg 自定义返回错误信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  code.getMsg(),
		"data": data,
	})
}
