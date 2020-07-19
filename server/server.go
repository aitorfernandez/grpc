package server

import (
	"context"

	"github.com/aitorfernandez/grpc/server/option"
	"google.golang.org/grpc"
)

// Server gRPC abstraction.
type Server interface {
	// Returns the Options.
	Options() *option.Options
	// Returns a gRPC server.
	GRPC() *grpc.Server
	// Start the server.
	Start() error
	// Stop the server.
	Stop()
	// UnaryInterceptor for RPCs calls.
	UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// New returns a new gRPC Server.
func New(opts ...option.Option) Server {
	return newGRPCServer(opts...)
}

// Address for manage option.Options.Address.
func Address(a string) option.Option {
	return func(o *option.Options) {
		o.Address = a
	}
}

// Before functions are executed on the gRPC object before the request.
func Before(before ...option.RequestFunc) option.Option {
	return func(o *option.Options) {
		o.Before = append(o.Before, before...)
	}
}

// After functions are executed on the gRPC object after the response.
func After(after ...option.ResponseFunc) option.Option {
	return func(o *option.Options) {
		o.After = append(o.After, after...)
	}
}
