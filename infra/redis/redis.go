/**
 * @Time: 2021/2/27 6:30 下午
 * @Author: varluffy
 */

package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func New(addr, password string, db int) (*redis.Client, func(), error) {
	opt := &redis.Options{
		Addr: addr,
		DB:   db,
	}
	if password != "" {
		opt.Password = password
	}
	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, nil, err
	}
	cleanFunc := func() {
		client.Close()
	}
	return client, cleanFunc, nil
}
