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
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/varluffy/rich/app"
	v1 "github.com/varluffy/rich/example/api/v1"
	"github.com/varluffy/rich/log"
	"github.com/varluffy/rich/transport/http"
	"github.com/varluffy/rich/transport/http/gin/middleware/logging"
	"github.com/varluffy/rich/transport/http/gin/middleware/recovery"
	"github.com/varluffy/rich/transport/http/gin/middleware/translation"
	"go.uber.org/zap"
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
		conf := viper.GetViper()
		logger := log.NewLogger(
			log.WithFileRotation(conf.GetString("log.out_put_file")),
			log.WithMaxAge(conf.GetInt("log.max_age")),
			log.WithMaxSize(conf.GetInt("log.max_size")),
			log.WithMaxBackups(conf.GetInt("log.max_backup")),
		)
		logger.Info("init app")
		app, clean, err := initApp(conf, logger)
		if err != nil {
			return err
		}
		defer func() {
			clean()
		}()
		return app.Run()
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

func newHttpServer(conf *viper.Viper, logger *zap.Logger) *http.Server {
	router := gin.New()
	gin.SetMode(conf.GetString("HTTP.Mode"))
	router.Use(
		recovery.Recovery(recovery.WithLogger(logger)),
		translation.Translation(),
		logging.Server(logging.WithLogger(logger)),
	)
	return http.NewServer(
		http.Address(viper.GetString("HTTP.addr")),
		http.Timeout(viper.GetDuration("HTTP.timeout")*time.Second),
		http.Logger(logger),
		http.Router(router),
	)
}

func newApp(logger *zap.Logger, hs *http.Server, router *v1.Router) *app.App {
	router.Register()
	return app.New(
		app.Name(Name),
		app.Version(Version),
		app.Logger(logger),
		app.Server(hs),
	)
}
