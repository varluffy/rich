// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/varluffy/rich/app"
	"github.com/varluffy/rich/example/api/middleware"
	v1 "github.com/varluffy/rich/example/api/v1"
	"github.com/varluffy/rich/example/internal/biz"
	"github.com/varluffy/rich/example/internal/repo"
	"github.com/varluffy/rich/example/internal/service"
	"go.uber.org/zap"
)

func initApp(conf *viper.Viper, logger *zap.Logger) (*app.App, func(), error) {
	panic(wire.Build(repo.ProviderSet, biz.ProviderSet, service.ProviderSet, middleware.ProviderSet, v1.NewRouter, newHttpServer, newApp))
}
