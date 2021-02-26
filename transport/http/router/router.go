/**
 * @Time: 2021/2/25 2:19 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/varluffy/ginx/log"
	"github.com/varluffy/ginx/transport/http/router/middleware"
	"go.uber.org/zap"
	"net/http"
)

func NewRouter(opts ...Option) *gin.Engine {
	opt := &option{
		logger:        log.NewLogger(),
		enableCors:    true,
		enableLogging: true,
		enablePProf:   false,
		panicNotify:   false,
		withoutLoggingPaths: make(map[string]bool),
	}
	for _, f := range opts {
		f(opt)
	}
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	if opt.enableLogging {
		engine.Use(middleware.Logging(opt.logger))
	}

	if opt.enablePProf {
		pprof.Register(engine)
	}

	if opt.enableCors {
		engine.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:     []string{"*"},
			AllowCredentials:   true,
			OptionsPassthrough: true,
		}))
	}
	engine.Use(middleware.Translation())
	return engine
}

func DisableLogging() Option {
	return func(opt *option) {
		opt.enableLogging = false
	}
}

func EnablePProf() Option {
	return func(opt *option) {
		opt.enablePProf = true
	}
}

func WithoutLoggingPaths(paths ...string) Option {
	return func(opt *option) {
		m := make(map[string]bool)
		for _, path := range paths {
			m[path] = true
		}
		opt.withoutLoggingPaths = m
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opt *option) {
		opt.logger = logger
	}
}

type Option func(opt *option)

type option struct {
	logger        *zap.Logger
	enableCors    bool
	enableLogging bool
	enablePProf   bool
	panicNotify   bool
	withoutLoggingPaths map[string]bool
}
