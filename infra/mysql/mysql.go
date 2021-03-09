/**
 * @Time: 2021/2/26 6:05 下午
 * @Author: varluffy
 */

package mysql

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Option func(option *options)

type options struct {
	logger       *zap.Logger
	maxLifeTime  time.Duration
	maxIdleConns int
	maxOpenConns int
	enableLog    bool
}

func New(dsn string, opts ...Option) (*gorm.DB, func(), error) {
	opt := &options{
		maxLifeTime:  time.Second * 7200,
		maxIdleConns: 150,
		maxOpenConns: 50,
		enableLog:    true,
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, nil, err
	}
	for _, o := range opts {
		o(opt)
	}
	sqlDB.SetConnMaxLifetime(opt.maxLifeTime)
	sqlDB.SetMaxIdleConns(opt.maxIdleConns)
	sqlDB.SetMaxOpenConns(opt.maxOpenConns)
	cleanFunc := func() {
		if err := sqlDB.Close(); err != nil {
			opt.logger.Error("close db error", zap.Error(err))
		}
	}
	return db, cleanFunc, nil
}

func WithDBMaxLifetime(duration time.Duration) Option {
	return func(opt *options) {
		opt.maxLifeTime = duration
	}
}

func WithDBMaxIdleConns(i int) Option {
	return func(opt *options) {
		opt.maxIdleConns = i
	}
}

func WithDBMaxOpenConns(o int) Option {
	return func(opt *options) {
		opt.maxOpenConns = o
	}
}
