/**
 * @Time: 2021/3/11 4:32 下午
 * @Author: varluffy
 */

package service

import (
	"context"
	v1 "github.com/varluffy/rich/example/api/app/v1"
	"github.com/varluffy/rich/example/internal/domain"
)

var _ v1.BlogServiceHTTPServer = &Article{}

type Article struct {
	usecase *domain.ArticleUseCase
}

func NewArticleService(usecase *domain.ArticleUseCase) *Article {
	return &Article{usecase: usecase}
}

func (a *Article) CreateArticle(ctx context.Context, article *v1.Article) (*v1.Article, error) {
	err := a.usecase.CreateArticle(ctx, &domain.Article{
		Title:   article.Title,
		Content: article.Content,
	})
	return article, err
}

func (a *Article) GetArticles(ctx context.Context, req *v1.GetArticlesReq) (*v1.GetArticlesResp, error) {
	return &v1.GetArticlesResp{
		Total:    1,
		Articles: nil,
	},nil
}