/**
 * @Time: 2021/3/11 3:36 下午
 * @Author: varluffy
 */

package server

import (
	"github.com/google/wire"
	"github.com/varluffy/rich/example/internal/domain"
	"github.com/varluffy/rich/example/internal/server/repo"
	"github.com/varluffy/rich/example/internal/server/service"
)

var Set = wire.NewSet(
	repo.NewArticleRepo,
	domain.NewArticleUsecase,
	service.NewArticleService,
)