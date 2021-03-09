/**
 * @Time: 2021/3/8 4:16 下午
 * @Author: varluffy
 */

package biz

import (
	"context"
	"go.uber.org/zap"
)

const (
	IdentityTypeMobile = iota + 1
	IdentityTypeWechat
)

type Auth struct {
	Id           int64
	UserId       int64
	IdentityType int64
	Identifier   string
	Unionid      string
	Credential   string
}

//IAuthRepository IAuthRepository
type AuthRepo interface {
	// db
	GetAuth(ctx context.Context, identityType int64, identifier string) (auth *Auth, err error)
	CreateAuth(ctx context.Context, auth *Auth) (err error)
	UpdateAuth(ctx context.Context, id int64, aut *Auth) (err error)
}

type AuthUsecase struct {
	repo   AuthRepo
	logger *zap.Logger
}

func NewAuthUsecase(repo AuthRepo, logger *zap.Logger) *AuthUsecase {
	return &AuthUsecase{repo: repo, logger: logger}
}

func (a *AuthUsecase) GetAuth(ctx context.Context, identityType int64, identifier string) (auth *Auth, err error) {
	return a.repo.GetAuth(ctx, identityType, identifier)
}
