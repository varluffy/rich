// +build wireinject

/**
 * @Time: 2021/3/11 3:33 下午
 * @Author: varluffy
 */

package main

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/varluffy/rich"
	"github.com/varluffy/rich/example/internal/server"
)

func newApp(conf *viper.Viper) (*rich.App, func(), error) {
	panic(wire.Build(server.Set, set))
}
