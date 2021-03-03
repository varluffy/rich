/**
 * @Time: 2021/3/2 11:16 上午
 * @Author: varluffy
 * @Description: service
 */

package service

import (
	"github.com/varluffy/ginx/example/api/middleware"
	"github.com/varluffy/ginx/example/internal/usecase"
	"go.uber.org/zap"
)

type Service struct {
	Auth *Auth
}

func NewService(logger *zap.Logger, u *usecase.Usecase, token middleware.Token) *Service {
	return &Service{
		Auth: NewAuth(logger, token, u.Auth),
	}
}
