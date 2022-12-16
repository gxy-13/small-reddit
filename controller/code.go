package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeSignUpFailed
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeWrongAuth
	CodeInvalidToken
)

var CodeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeSignUpFailed:    "注册失败",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "请登录",
	CodeWrongAuth:       "请求头中auth格式有误",
	CodeInvalidToken:    "无效的token",
}

func (code ResCode) getMsg() string {
	msg, ok := CodeMsgMap[code]
	if !ok {
		return CodeMsgMap[CodeServerBusy]
	}
	return msg
}
