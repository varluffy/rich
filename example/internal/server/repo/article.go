/**
 * @Time: 2021/3/11 4:29 下午
 * @Author: varluffy
 */

package repo

import (
	"context"
	"github.com/varluffy/rich/example/internal/domain"
	"gorm.io/gorm"
)

type article struct {
	db *gorm.DB
}

func (a article) GetArticle(ctx context.Context, id int) (*domain.Article, error) {
	return &domain.Article{
		Title:   "123",
		Content: "23123",
	}, nil
}

func (a article) CreateArticle(ctx context.Context, article *domain.Article) error {
	return nil
}

func NewArticleRepo(db *gorm.DB) domain.IArticleRepo {
	return &article{db: db}
}

