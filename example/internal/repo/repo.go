/**
 * @Time: 2021/2/26 5:00 下午
 * @Author: varluffy
 */

package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"time"
)

var ProviderSet = wire.NewSet(NewRepo, NewAuthRepo)

type Repo struct {
	db     *gorm.DB
	rds    *redis.Client
	logger *zap.Logger
}

func NewRepo(conf *viper.Viper, logger *zap.Logger) (*Repo, func(), error) {
	logger = logger.With(zap.String("module", "repo"))

	// initialize mysql
	db, err := gorm.Open(mysql.Open(conf.GetString("mysql.dsn")), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Info),
	})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	if d := conf.GetDuration("mysql.max_life_time"); d > 0 {
		sqlDB.SetConnMaxLifetime(d * time.Second)
	}

	if n := conf.GetInt("mysql.max_open_conns"); n > 0 {
		sqlDB.SetMaxOpenConns(n)
	}

	if n := conf.GetInt("mysql.max_idle_conns"); n > 0 {
		sqlDB.SetMaxIdleConns(n)
	}

	if t := conf.GetBool("mysql.auto_migrate"); t {
		if err := db.AutoMigrate(); err != nil {
			return nil, nil, err
		}
	}

	// initialize redis
	rdsOption := &redis.Options{
		Addr: conf.GetString("redis.addr"),
		DB:   0,
	}
	if pass := conf.GetString("redis.pass"); pass != "" {
		rdsOption.Password = pass
	}

	if db := conf.GetInt("redis.db"); db > 0 {
		rdsOption.DB = db
	}

	rds := redis.NewClient(rdsOption)
	if err := rds.Ping(context.Background()).Err(); err != nil {
		return nil, nil, err
	}
	return &Repo{
			db:     db,
			rds:    rds,
			logger: logger,
		}, func() {
			if err := sqlDB.Close(); err != nil {
				logger.Error("db.close error", zap.Error(err))
			}
			if err := rds.Close(); err != nil {
				logger.Error("redis.close error", zap.Error(err))
			}
			logger.Info("close repo ...")
		}, nil

}
