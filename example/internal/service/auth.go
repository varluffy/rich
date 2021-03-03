/**
 * @Time: 2021/3/2 11:18 上午
 * @Author: varluffy
 * @Description: auth user service
 */

package service

import (
	"github.com/gin-gonic/gin"
	"github.com/varluffy/ginx/example/api/middleware"
	"github.com/varluffy/ginx/example/internal/usecase"
	"github.com/varluffy/ginx/transport/http/router/ginwrap"
	"go.uber.org/zap"
)

type Auth struct {
	usecase *usecase.Auth
	token   middleware.Token
	logger  *zap.Logger
}

func NewAuth(logger *zap.Logger, token middleware.Token, usecase *usecase.Auth) *Auth {
	return &Auth{
		usecase: usecase,
		token:   token,
		logger:  logger,
	}
}

type LoginRequest struct {
	IdentityType int64  `form:"identityType" binding:"required"`
	Identifier   string `form:"identifier" binding:"required"`
	SmsCode      string `form:"smsCode" binding:"required"`
}

func (a *Auth) Login(c *gin.Context) {
	var req LoginRequest
	if err := ginwrap.BindAndValid(c, &req); err != nil {
		ginwrap.ErrorResponse(c, err)
		return
	}
	userOauth, err := a.usecase.GetAuth(c.Request.Context(), req.IdentityType, req.Identifier)
	if err != nil {
		a.logger.Error("get auth error", zap.Error(err))
		ginwrap.ErrorResponse(c, err)
		return
	}
	token, _ := a.token.Sign(userOauth.UserId)
	ginwrap.Response(c, gin.H{"token": token})
}
