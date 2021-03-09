/**
 * @Time: 2021/3/4 2:41 下午
 * @Author: varluffy
 */

package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/varluffy/rich/example/api/code"
	"github.com/varluffy/rich/transport/http/gin/ginx"
	"strings"
)

const (
	userIdKey = "auth-userId"
)

var ProviderSet = wire.NewSet(NewToken, NewMiddleware)

type Middleware struct {
	token Token
}

func NewMiddleware(token Token) *Middleware {
	return &Middleware{
		token: token,
	}
}

func (a *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ginx.ErrorResponse(c, code.ErrUnauthorizedInvalid)
			return
		}
		tokenString = tokenString[7:]
		claim, err := a.token.Parse(tokenString)
		if err != nil {
			switch err.(jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				ginx.Response(c, code.ErrUnauthorizedExpired)
				return
			default:
				ginx.ErrorResponse(c, code.ErrUnauthorizedError)
				return
			}
		}
		if claim.UserId <= 0 {
			ginx.ErrorResponse(c, code.ErrUnauthorizedInvalid)
			return
		}
		c.Set(userIdKey, claim.UserId)
		c.Next()
	}
}

func FromUserId(c *gin.Context) int64 {
	return c.GetInt64(userIdKey)
}
