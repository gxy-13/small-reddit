package middleware

import (
	"awesomeProject/controller"
	"awesomeProject/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JwtAuth JWT的认证中间件
func JwtAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token的三种方式，1.请求头 2.请求体 3. 放在URI
		// 假设Token放在Header的Authorization中，使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		// Token不满足格式要求
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeWrongAuth)
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将解析token后得到的username保存在请求的上下文中,其余函数可以通过c.Get("username")来获取当前请求的用户信息
		c.Set(controller.ContextUserIDKey, mc.UserID)
		c.Next()
	}
}
