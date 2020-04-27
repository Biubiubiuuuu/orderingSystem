package jwtHelper

import (
	"time"

	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(configHelper.JwtSecret)

type Claims struct {
	UserName string
	Password string
	jwt.StandardClaims
}

// JWT加密生成并返回令牌
// param username
// param password
// return string,error
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    configHelper.JwtName,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// JWT解析令牌并验证
// param tokenStr
// retrun *Claims, error
func ParseToken(tokenStr string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
