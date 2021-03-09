/**
 * @Time: 2021/2/27 7:03 下午
 * @Author: varluffy
 */

package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type Token interface {
	Sign(userId int64) (token string, err error)
	Parse(token string) (*claims, error)
}

type claims struct {
	UserId int64
	jwt.StandardClaims
}

type token struct {
	secret string
	expire time.Duration
}

func NewToken(conf *viper.Viper) Token {
	return &token{secret: conf.GetString("jwt.secret"), expire: conf.GetDuration("jwt.expire")}
}

func (t *token) Sign(userId int64) (tokenString string, err error) {
	claims := claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(t.expire * time.Hour).Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.secret))
}

func (t *token) Parse(tokenString string) (c *claims, err error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
