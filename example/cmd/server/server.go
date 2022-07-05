/**
 * @Time: 2021/3/11 3:33 下午
 * @Author: varluffy
 */

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/varluffy/rich"
	v1 "github.com/varluffy/rich/example/api/app/v1"
	"github.com/varluffy/rich/example/internal/server/service"
	"github.com/varluffy/rich/log"
	"github.com/varluffy/rich/transport/grpc"
	"github.com/varluffy/rich/transport/http"
	"github.com/varluffy/rich/transport/http/gin/ginx"
	"github.com/varluffy/rich/transport/http/gin/middleware/logging"
	"github.com/varluffy/rich/transport/http/gin/middleware/recovery"
	"github.com/varluffy/rich/transport/http/gin/middleware/translation"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"time"
)

var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// cfgFile config file
	cfgFile string
)
var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "example",
	Long:  "example",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := newConfig()
		if err != nil {
			return err
		}
		app, clean, err := newApp(conf)
		if err != nil {
			return err
		}
		defer clean()
		return app.Run()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "../../config/config.yaml", "config file")
}

func newConfig() (v *viper.Viper, err error) {
	v = viper.New()
	v.SetConfigFile(cfgFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return
}

var set = wire.NewSet(initLogger, initHttpServer, initApp, initMysql, initRedis, initGrpcServer, wire.Struct(new(services), "*"))

func initLogger(conf *viper.Viper) *zap.Logger {
	logger := log.NewLogger(
		log.WithFileRotation(conf.GetString("log.out_put_file")),
		log.WithMaxAge(conf.GetInt("log.max_age")),
		log.WithMaxSize(conf.GetInt("log.max_size")),
		log.WithMaxBackups(conf.GetInt("log.max_backup")),
	)
	logger.Info("init app")
	return logger
}

func initMysql(conf *viper.Viper, logger *zap.Logger) (*gorm.DB, func(), error) {
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
	return db, func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("db.close error", zap.Error(err))
		}
		logger.Info("close mysql ...")
	}, nil
}

func initRedis(conf *viper.Viper, logger *zap.Logger) (*redis.Client, func(), error) {
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
	return rds, func() {
		if err := rds.Close(); err != nil {
			logger.Error("redis.close error", zap.Error(err))
		}
		logger.Info("close repo ...")
	}, nil
}

func initGrpcServer(conf *viper.Viper, logger *zap.Logger, svc *services) (*grpc.Server, error) {
	gs := grpc.NewServer(
		grpc.Address(conf.GetString("grpc.addr")),
	)
	logger.Info("init grpc server ...", zap.String("addr", conf.GetString("grpc.addr")))
	v1.RegisterBlogServiceServer(gs, svc.article)
	return gs, nil
}

func initHttpServer(conf *viper.Viper, logger *zap.Logger, svc *services) (*http.Server, error) {
	if err := ginx.TransInit("zh"); err != nil {
		return nil, err
	}
	gin.SetMode(conf.GetString("http.mode"))
	e := gin.New()
	e.Use(
		recovery.Recovery(recovery.WithLogger(logger)),
		translation.Translation(),
		logging.Server(logging.WithLogger(logger)),
	)
	svc.register(e)
	hs := http.NewServer(
		http.Address(conf.GetString("http.addr")),
		http.Timeout(conf.GetDuration("http.timeout")*time.Second),
		http.Logger(logger),
		http.Router(e),
	)
	return hs, nil
}

func initApp(logger *zap.Logger, hs *http.Server, gs *grpc.Server) *rich.App {
	return rich.New(
		rich.Name(Name),
		rich.Version(Version),
		rich.Logger(logger),
		rich.Server(hs, gs),
	)
}

type services struct {
	article *service.Article
}

// 添加路由
func (s *services) register(r gin.IRouter) {
	v1.RegisterBlogServiceHTTPServer(r, s.article)
}
