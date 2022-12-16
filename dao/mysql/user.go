package mysql

import (
	"awesomeProject/model"
	"awesomeProject/utils/encrypt"
	sql2 "database/sql"
	errors "errors"

	"go.uber.org/zap"
)

// 用户数据操作

var (
	ErrorUserExist    = errors.New("用户已存在")
	ErrorUserNotExist = errors.New("用户名不存在")
	ErrorPassword     = errors.New("密码错误")
)

// CheckUserExist 判断用户是否存在
func CheckUserExist(username string) (err error) {
	// 使用count(user_id) 来判断是否存在用户
	sql := `select count(user_id) from user where username = ?`
	var count int
	err = db.Get(&count, sql, username)
	// 数据库查询失败
	if err != nil {
		zap.L().Error("mysql.CheckUserExist failed", zap.Error(err))
		return err
	}
	if count > 0 {
		zap.L().Error("用户已存在", zap.String("username", username), zap.Error(err))
		return ErrorUserExist
	}
	return
}

// InsertUser 注册用户
func InsertUser(p *model.User) (err error) {
	// 需要先对password进行加密
	p.Password = encrypt.CryptPassword(p.Password)
	sql := `insert into user(user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sql, p.UserID, p.Username, p.Password)
	return err
}

// Login 用户登陆
func Login(u *model.User) (err error) {
	// 需要先对密码进行加密
	CPassword := encrypt.CryptPassword(u.Password)
	sql := `select user_id, username, password from user where username = ?`
	err = db.Get(u, sql, u.Username)
	// 没有数据
	if err == sql2.ErrNoRows {
		zap.L().Error("用户名不存在", zap.String("username", u.Username), zap.Error(err))
		return ErrorUserNotExist
	}
	// 查询失败
	if err != nil {
		zap.L().Error("mysql.Login() failed", zap.Error(err))
		return
	}
	// 密码错误
	if CPassword != u.Password {
		zap.L().Error("密码错误", zap.Error(err))
		return ErrorPassword
	}
	return
}

// GetUserByID 通过uid 查找用户信息
func GetUserByID(uid int64) (username string, err error) {
	sql := `select username from user where user_id = ?`
	err = db.Get(&username, sql, uid)
	if err != nil {
		return
	}
	return
}
