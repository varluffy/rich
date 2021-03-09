/**
 * @Time: 2021/3/2 10:36 上午
 * @Author: varluffy
 */

package repo

import (
	"context"
	"github.com/varluffy/rich/example/internal/biz"
	"github.com/varluffy/rich/example/internal/repo/entity"
	"go.uber.org/zap"
)

var _ biz.AuthRepo = (*UserAuth)(nil)

type UserAuth struct {
	repo   *Repo
	logger *zap.Logger
}

func NewAuthRepo(repo *Repo, logger *zap.Logger) biz.AuthRepo {
	return &UserAuth{
		repo:   repo,
		logger: logger,
	}
}

func (u UserAuth) GetAuth(ctx context.Context, identityType int64, identifier string) (auth *biz.Auth, err error) {
	var a entity.UserOauths
	err = u.repo.db.Where("identity_type = ? and identifier = ?", identityType, identifier).First(&a).Error
	if err != nil {
		return nil, err
	}
	auth = &biz.Auth{
		Id:           a.ID,
		UserId:       a.UserId,
		IdentityType: a.IdentityType,
		Identifier:   a.Identifier,
		Unionid:      a.UnionId,
		Credential:   a.Credential,
	}
	return
}

func (u UserAuth) CreateAuth(ctx context.Context, auth *biz.Auth) (err error) {
	panic("implement me")
}

func (u UserAuth) UpdateAuth(ctx context.Context, id int64, aut *biz.Auth) (err error) {
	panic("implement me")
}
