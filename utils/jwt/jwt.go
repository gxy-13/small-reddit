package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims 自定义声明类型，并内嵌jwt.RegisteredClaims
// 自带的jwt.RegisteredClaims只包含了默认的7个字段

type CustomClaims struct {
	// 根据需要添加字段
	UserID               int64  `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 默认的7个标准声明
}

// TokenExpireDuration JWT过期时间
const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于签名的字符串
var CustomSecret = []byte("Fire")

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	//创建自己的声明
	claims := CustomClaims{
		userID,
		username,
		jwt.RegisteredClaims{
			Issuer:    "FREE",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
		},
	}
	// 使用指定签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 如果是标准的Claims则可以直接使用Parse
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token)....)
	// 如果是自定义Claims结构体则需要使用 parseWithClaims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象的Claims进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
