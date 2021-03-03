/**
 * @Time: 2021/2/27 7:19 下午
 * @Author: varluffy
 * @Description: auth
 */

package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/varluffy/ginx/example/api/code"
	"github.com/varluffy/ginx/transport/http/router/ginwrap"
	"strings"
)

const (
	userIdKey = "auth-userId"
)

type AuthMiddleware struct {
	token Token
}

func NewAuth(token Token) *AuthMiddleware {
	return &AuthMiddleware{token: token}
}

func (a *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ginwrap.ErrorResponse(c, code.ErrUnauthorizedInvalid)
			return
		}
		tokenString = tokenString[7:]
		claim, err := a.token.Parse(tokenString)
		if err != nil {
			switch err.(jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				ginwrap.Response(c, code.ErrUnauthorizedExpired)
				return
			default:
				ginwrap.ErrorResponse(c, code.ErrUnauthorizedError)
				return
			}
		}
		if claim.UserId < 0 {
			ginwrap.ErrorResponse(c, code.ErrUnauthorizedInvalid)
			return
		}
		c.Set(userIdKey, claim.UserId)
		c.Next()
	}
}

func FromUserId(c *gin.Context) int64 {
	return c.GetInt64(userIdKey)
}
