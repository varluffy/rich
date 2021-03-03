/**
 * @Time: 2021/3/2 10:36 上午
 * @Author: varluffy
 * @Description: userauth
 */

package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/varluffy/ginx/example/internal/domain"
	"github.com/varluffy/ginx/example/internal/repo/entity"
	"gorm.io/gorm"
)

var _ domain.AuthRepo = (*UserAuth)(nil)

type UserAuth struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewUserAuthRepo(db *gorm.DB, rds *redis.Client) domain.AuthRepo {
	return &UserAuth{
		db:  db,
		rds: rds,
	}
}

func (u UserAuth) GetAuth(ctx context.Context, identityType int64, identifier string) (auth *domain.Auth, err error) {
	var a entity.UserOauths
	err = u.db.Where("identity_type = ? and identifier = ?", identityType, identifier).First(&a).Error
	if err != nil {
		return nil, err
	}
	auth = &domain.Auth{
		Id:           a.ID,
		UserId:       a.UserId,
		IdentityType: a.IdentityType,
		Identifier:   a.Identifier,
		Unionid:      a.UnionId,
		Credential:   a.Credential,
	}
	return
}

func (u UserAuth) CreateAuth(ctx context.Context, auth *domain.Auth) (err error) {
	panic("implement me")
}

func (u UserAuth) UpdateAuth(ctx context.Context, id int64, aut *domain.Auth) (err error) {
	panic("implement me")
}

func (u UserAuth) SetCode(ctx context.Context, mobile string, code string) (err error) {
	panic("implement me")
}

func (u UserAuth) CheckCode(ctx context.Context, mobile string, code string) (err error) {
	panic("implement me")
}
