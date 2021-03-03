/**
 * @Time: 2021/2/26 5:00 下午
 * @Author: varluffy
 * @Description: repo
 */

package repo

import (
	"github.com/go-redis/redis/v8"
	"github.com/varluffy/ginx/example/internal/domain"
	"gorm.io/gorm"
)

type Repo struct {
	AuthRepo domain.AuthRepo
}

func NewRepo(db *gorm.DB, Rds *redis.Client) *Repo {
	return &Repo{
		AuthRepo: NewUserAuthRepo(db, Rds),
	}
}
