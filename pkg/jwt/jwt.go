package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// 定义Token过期时间
const AccessTokenExpireDuration = time.Hour * 24
const RefreshTokenExpireDuration = time.Hour * 24 * 7

var mySecret = []byte("Mitsui515")

func keyFunc(_ *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

// MyClaims 自定义的 metadata 在加密后作为 JWT 的第二部分返回给客户端
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成 JWT
func GenToken(userID int64, username string) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("auth.jwt_atoken_expire") * time.Hour).Unix(), // 过期时间
			Issuer:    "minireddit",                                                                   // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// refresh token 不需要存任何自定义信息，所以不需要自定义的 MyClaims
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(viper.GetDuration("auth.jwt_rtoken_expire") * time.Hour).Unix(), // 过期时间
		Issuer:    "minireddit",                                                                   // 签发人
	}).SignedString(mySecret)
	return
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	// 校验token
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken 刷新 JWT
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token 无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧 access token 中解析出 claims 数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	// 当 access token 是过期错误 并且 refresh token 未过期，重新生成 access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}
