/**
 * @Time: 2021/2/24 8:46 下午
 * @Author: varluffy
 * @Description: //app
 */

package app

import (
	"context"
	"errors"
	"github.com/varluffy/ginx/log"
	"github.com/varluffy/ginx/transport"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	opt    option
	ctx    context.Context
	cancel func()
	logger *zap.Logger
}

func New(opts ...Option) *App {
	opt := option{
		ctx:     context.Background(),
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		logger:  log.NewLogger(),
		servers: nil,
	}
	for _, o := range opts {
		o(&opt)
	}
	ctx, cancel := context.WithCancel(opt.ctx)
	return &App{
		opt:    opt,
		ctx:    ctx,
		cancel: cancel,
		logger: opt.logger,
	}
}

func (a *App) Logger() *zap.Logger {
	return a.opt.logger
}

func (a *App) Server() []transport.Server {
	return a.opt.servers
}

func (a *App) Run() error {
	a.logger.Info("app run", zap.String("service_name", a.opt.name), zap.String("version", a.opt.version))
	g, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.opt.servers {
		server := srv
		g.Go(func() error {
			<-ctx.Done()
			return server.Stop()
		})
		g.Go(func() error {
			return server.Start()
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opt.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return a.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

type Option func(opt *option)

type option struct {
	name    string
	version string

	ctx  context.Context
	sigs []os.Signal

	logger  *zap.Logger
	servers []transport.Server
}

func Name(name string) Option {
	return func(opt *option) {
		opt.name = name
	}
}

func Version(version string) Option {
	return func(opt *option) {
		opt.version = version
	}
}

func Server(srv ...transport.Server) Option {
	return func(opt *option) {
		opt.servers = srv
	}
}

func Logger(logger *zap.Logger) Option {
	return func(opt *option) {
		opt.logger = logger
	}
}
