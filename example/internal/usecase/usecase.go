/**
 * @Time: 2021/3/2 3:40 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package usecase

import (
	"github.com/varluffy/ginx/example/internal/repo"
	"go.uber.org/zap"
)

type Usecase struct {
	Auth *Auth
}

func NewUsecase(repo *repo.Repo, logger *zap.Logger) *Usecase {
	return &Usecase{
		Auth: NewAuth(repo.AuthRepo, logger),
	}
}
