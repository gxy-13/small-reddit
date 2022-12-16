package logic

import (
	"awesomeProject/dao/mysql"
	"awesomeProject/model"
	"awesomeProject/utils/jwt"
	"awesomeProject/utils/snowflake"

	"go.uber.org/zap"
)

// SignUp 注册的整体逻辑，先判断输入的用户名是否存在，再生成UID，最后注册用户信息
func SignUp(p *model.SignUpParam) (err error) {
	// 判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		zap.L().Error("user is exist")
		return err
	}
	// 生成UID
	userID := snowflake.GenID()
	// 生成一个User对象
	var u = model.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 将用户信息插入数据库
	return mysql.InsertUser(&u)
}

// Login 用户登陆逻辑，直接调用dao层查询数据库
func Login(u *model.LoginUser) (user *model.User, err error) {
	user = &model.User{
		Password: u.Password,
		Username: u.Username,
	}
	err = mysql.Login(user)
	if err != nil {
		return
	}
	// 生成token并返回给controller
	tokenString, err := jwt.GenToken(user.UserID, user.Username)
	user.Token = tokenString
	return
}
