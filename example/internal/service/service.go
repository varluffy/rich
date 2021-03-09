/**
 * @Time: 2021/3/8 4:29 下午
 * @Author: varluffy
 */

package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewService, NewAuthService)

type Service struct {
	Auth *AuthService
}

func NewService(auth *AuthService) *Service {
	return &Service{
		Auth: auth,
	}
}
