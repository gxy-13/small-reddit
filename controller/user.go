package controller

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/logic"
	"awesomeProject/model"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// ContextUserIDKey 上下文中存放的UserID， 将代码钉死的字符串定义成一个常量
const ContextUserIDKey = "userID"

var (
	ErrorNoLogin = errors.New("请登录")
)

func SignUp(c *gin.Context) {
	// 参数获取和参数交验
	// 创建结构体指针，提升性能
	p := new(model.SignUpParam)
	// 前端传送的是JSON格式，使用ShouldBindJSON 用结构体来绑定数据
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Sign up with invalid param")
		// 首先判断报错类型是否为validator.ValidationErrors类型的错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("invalid param")
			// 如果不是校验器类型错误直接返回
			ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
			return
		}
		// 是校验器错误类型就进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}
	//// 手动校验参数是否满足业务逻辑
	//if len(p.Password) == 0 || len(p.Username) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
	//	zap.L().Error("Sign up with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求的参数有误",
	//	})
	//	return
	//}
	// 业务逻辑
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("user sign up failed")
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	zap.L().Info("user sign up success")
	ResponseSuccess(c, nil)

}

// Login 用户登陆
func Login(c *gin.Context) {
	// 获取参数以及参数校验
	// 获取一个LoginUser实例,创建结构体指针提升性能
	u := new(model.LoginUser)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("Login with invalid param")
		// 判断是否为校验器错误
		errs, ok := err.(validator.ValidationErrors)
		// 不是校验器错误就直接返回
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 返回翻译信息
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}
	// 登陆逻辑
	// 创建一个User实例，存放登陆信息
	//user := new(model.User)
	user, err := logic.Login(u)
	if err != nil {
		zap.L().Error("user login failed")
		if errors.Is(err, mysql.ErrorPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回登陆成功
	zap.L().Info("user login success")
	ResponseSuccess(c, gin.H{
		// 将uid转换为string给前端，json的大小是 -2<<53 -1 ~ 2 << 53 - 1
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}

// GetCurrentUserID 通过存放在上下文中的ContextUserIDKey 获取当前userID
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorNoLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorNoLogin
		return
	}
	return
}
