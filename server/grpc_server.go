package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aitorfernandez/grpc/server/option"
	"google.golang.org/grpc"
)

type grpcServer struct {
	Srv     *grpc.Server
	Opts    *option.Options
	running bool
}

// setup options and creates a gRPC server which has no service registered.
func (s *grpcServer) setup() {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(s.UnaryInterceptor),
	}

	s.Srv = grpc.NewServer(opts...)
}

// UnaryInterceptor intercepts the unary RPCs and execute callbacks if exists.
func (s *grpcServer) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	c := s.Options()
	for _, f := range c.Before {
		f(ctx, req)
	}

	res, err := handler(ctx, req)
	for _, f := range c.After {
		f(ctx, res)
	}

	return res, err
}

// Options returns the grpcServer options.
func (s *grpcServer) Options() *option.Options {
	return s.Opts
}

// GRPC returns the gRPC Server.
func (s *grpcServer) GRPC() *grpc.Server {
	return s.Srv
}

// Start starts the server.
func (s *grpcServer) Start() error {
	if s.running {
		return nil
	}

	c := s.Options()
	lis, err := net.Listen("tcp", c.Address)
	if err != nil {
		return fmt.Errorf("grpc/server/Start net.Listen %w", err)
	}

	go func() {
		if err := s.GRPC().Serve(lis); err != nil {
			log.Fatal("grpc/server.Start Server start error %w", err)
		}
	}()

	s.running = true
	return nil
}

// Stop stops the server.
func (s *grpcServer) Stop() {
	if !s.running {
		return
	}
	s.GRPC().GracefulStop()
}

func newGRPCServer(opts ...option.Option) Server {
	o := option.New(opts...)
	srv := &grpcServer{
		Opts: o,
	}
	srv.setup()
	return srv
}
