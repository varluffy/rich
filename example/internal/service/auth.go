/**
 * @Time: 2021/3/8 4:29 下午
 * @Author: varluffy
 */

package service

import (
	"github.com/gin-gonic/gin"
	"github.com/varluffy/rich/example/api/middleware"
	"github.com/varluffy/rich/example/internal/biz"
	"github.com/varluffy/rich/transport/http/gin/ginx"
	"go.uber.org/zap"
)

type AuthService struct {
	logger *zap.Logger
	auth   *biz.AuthUsecase
	token  middleware.Token
}

func NewAuthService(
	auth *biz.AuthUsecase,
	logger *zap.Logger,
	token middleware.Token,
) *AuthService {
	return &AuthService{
		logger: logger,
		auth:   auth,
		token:  token,
	}
}

type LoginRequest struct {
	IdentityType int64  `form:"identityType" binding:"required"`
	Identifier   string `form:"identifier" binding:"required"`
}

func (s *AuthService) Login(c *gin.Context) {
	var req LoginRequest
	if err := ginx.ShouldBind(c, &req); err != nil {
		ginx.ErrorResponse(c, err)
		return
	}
	userOauth, err := s.auth.GetAuth(c.Request.Context(), req.IdentityType, req.Identifier)
	if err != nil {
		s.logger.Error("get auth error", zap.Error(err))
		ginx.ErrorResponse(c, err)
		return
	}
	token, _ := s.token.Sign(userOauth.UserId)
	ginx.Response(c, gin.H{"token": token})
}
