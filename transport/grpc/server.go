/**
 * @Time: 2022/5/19 09:36
 * @Author: varluffy
 */

package grpc

import (
	"context"
	"crypto/tls"
	"github.com/varluffy/rich/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"time"
)

var (
	_ transport.Server = (*Server)(nil)
)

type ServerOption func(o *Server)

func Network(network string) ServerOption {
	return func(o *Server) {
		o.network = network
	}
}

func Address(address string) ServerOption {
	return func(o *Server) {
		o.address = address
	}
}

func Timeout(timeout time.Duration) ServerOption {
	return func(o *Server) {
		o.timeout = timeout
	}
}

type Server struct {
	*grpc.Server
	baseCtx    context.Context
	network    string
	address    string
	timeout    time.Duration
	tlsConf    *tls.Config
	grpcOpts   []grpc.ServerOption
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	logger     *zap.Logger
	err        error
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx: context.Background(),
		network: "tcp",
		address: ":0",
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(srv)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(srv.unaryInts...),
		grpc.ChainStreamInterceptor(srv.streamInts...),
	}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	return srv
}

func (s *Server) Start() error {
	s.logger.Info("grpc server start", zap.String("address", s.address))
	go func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			s.logger.Error("grpc server start error", zap.Error(err))
			return
		}
		if s.err = s.Serve(lis); s.err != nil {
			s.logger.Error("grpc server start error", zap.Error(s.err))
			return
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("grpc server stop")
	s.Server.GracefulStop()
	return nil
}
