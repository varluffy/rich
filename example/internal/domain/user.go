/**
 * @Time: 2021/3/2 4:25 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package domain

import "context"

type User struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

type UserRepo interface {
	Info(ctx context.Context, userId int64) (user *User, err error)
}
