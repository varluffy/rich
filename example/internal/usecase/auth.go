/**
 * @Time: 2021/3/2 3:41 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package usecase

import (
	"context"
	"github.com/varluffy/ginx/example/internal/domain"
	"go.uber.org/zap"
)

type Auth struct {
	repo   domain.AuthRepo
	logger *zap.Logger
}

func NewAuth(repo domain.AuthRepo, logger *zap.Logger) *Auth {
	return &Auth{repo: repo, logger: logger}
}

func (au *Auth) GetAuth(ctx context.Context, identityType int64, identifier string) (auth *domain.Auth, err error) {
	return au.repo.GetAuth(ctx, identityType, identifier)
}

func (au *Auth) SmsCode(ctx context.Context, mobile string) error {
	//todo send smscode
	return au.repo.SetCode(ctx, mobile, "123456")
}

func (au *Auth) CheckCode(ctx context.Context, mobile, smsCode string) error {
	return au.repo.CheckCode(ctx, mobile, smsCode)
}
