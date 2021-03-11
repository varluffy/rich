/**
 * @Time: 2021/3/11 4:23 下午
 * @Author: varluffy
 */

package domain

import (
	"context"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title string `json:"title"`
	Content string `json:"content"`
}

type IArticleRepo interface {
	GetArticle(ctx context.Context, id int) (*Article, error)
	CreateArticle(ctx context.Context, article *Article) error
}

type ArticleUseCase struct {
	repo IArticleRepo
}

func NewArticleUsecase(repo IArticleRepo) *ArticleUseCase {
	return &ArticleUseCase{repo: repo}
}

func (u *ArticleUseCase) GetArticle(ctx context.Context, id int) (*Article, error) {
	return u.repo.GetArticle(ctx, id)
}

func (u *ArticleUseCase) CreateArticle(ctx context.Context, article *Article) error {
	return u.repo.CreateArticle(ctx, article)
}