/**
 * @Time: 2021/2/26 6:05 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package repository

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DBOption func(option *dbOption)

type dbOption struct {
	logger       *zap.Logger
	maxLifeTime  time.Duration
	maxIdleConns int
	maxOpenConns int
	enableLog    bool
}

func NewDB(dsn string, opts ...DBOption) (*gorm.DB, func(), error) {
	opt := &dbOption{
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

func WithDBMaxLifetime(duration time.Duration) DBOption {
	return func(opt *dbOption) {
		opt.maxLifeTime = duration
	}
}

func WithDBMaxIdleConns(i int) DBOption {
	return func(opt *dbOption) {
		opt.maxIdleConns = i
	}
}

func WithDBMaxOpenConns(o int) DBOption {
	return func(opt *dbOption) {
		opt.maxOpenConns = o
	}
}
