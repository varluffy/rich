/**
 * @Time: 2021/3/1 1:27 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package domain

import (
	"context"
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

type AuthRepo interface {
	// db
	GetAuth(ctx context.Context, identityType int64, identifier string) (auth *Auth, err error)
	CreateAuth(ctx context.Context, auth *Auth) (err error)
	UpdateAuth(ctx context.Context, id int64, aut *Auth) (err error)

	// redis
	SetCode(ctx context.Context, mobile string, code string) (err error)
	CheckCode(ctx context.Context, mobile string, code string) (err error)
}
