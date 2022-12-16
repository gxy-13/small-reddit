package model

// SignUpParam 注册时的必填信息
type SignUpParam struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// User 用户信息
type User struct {
	UserID   int64  `db:"user_id" json:"user_id,string"`
	Username string `db:"username" json:"user_name"`
	Password string `db:"password"`
	Token    string `json:"token"`
}

// LoginUser 登陆时的必填信息
type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
