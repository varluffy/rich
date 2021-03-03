/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/varluffy/ginx/app"
	"github.com/varluffy/ginx/example/api/middleware"
	v1 "github.com/varluffy/ginx/example/api/v1"
	"github.com/varluffy/ginx/example/internal/repo"
	"github.com/varluffy/ginx/example/internal/service"
	"github.com/varluffy/ginx/example/internal/usecase"
	"github.com/varluffy/ginx/log"
	"github.com/varluffy/ginx/repository"
	"github.com/varluffy/ginx/transport/http"
	"github.com/varluffy/ginx/transport/http/router"
	"os"
	"time"
)

var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "application",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		return initApp()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "../../config/config.yaml", "config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgFile)
	var err error
	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		return
	}
	panic("viper.ReadInConfig error:" + cfgFile + err.Error())
}

func initApp() error {
	logger := log.NewLogger(
		log.WithFileRotation(viper.GetString("log.OutpufFile")),
		log.WithMaxAge(viper.GetInt("log.MaxAge")),
		log.WithMaxSize(viper.GetInt("log.MaxSize")),
		log.WithMaxBackups(viper.GetInt("log.Backups")),
	)
	logger.Info("init app")
	db, cleanDBFunc, err := repository.NewDB(
		viper.GetString("MySQL.DSN"),
		repository.WithDBMaxLifetime(viper.GetDuration("MySQL.MaxLifetime")*time.Second),
		repository.WithDBMaxIdleConns(viper.GetInt("MySQL.MaxIdleConns")),
		repository.WithDBMaxOpenConns(viper.GetInt("MySQL.MaxOpenConns")),
	)
	if err != nil {
		return err
	}
	if viper.GetBool("MySQL.EnableAutoMigrate") {
		// todo add migrate
	}
	rds, cleanRdsFunc, err := repository.NewRedis(
		viper.GetString("Redis.Addr"),
		viper.GetString("Redis.Password"),
		viper.GetInt("Redis.DB"),
	)
	if err != nil {
		return err
	}
	defer func() {
		cleanDBFunc()
		cleanRdsFunc()
	}()
	engine := router.NewRouter(
		router.WithLogger(logger),
	)
	httpServer := http.NewServer(
		http.Address(viper.GetString("HTTP.addr")),
		http.Timeout(viper.GetDuration("HTTP.timeout")*time.Second),
		http.Logger(logger),
		http.Router(engine),
	)
	application := app.New(
		app.Name(Name),
		app.Version(Version),
		app.Logger(logger),
		app.Server(httpServer),
	)
	repo := repo.NewRepo(db, rds)
	u := usecase.NewUsecase(repo, logger)
	token := middleware.NewToken(viper.GetString("Token.Secret"), viper.GetDuration("Token.Expire"))
	svc := service.NewService(logger, u, token)
	v1.RegisterRouter(engine, svc, token)
	return application.Run()
}
