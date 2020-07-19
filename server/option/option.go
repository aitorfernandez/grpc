package option

import (
	"context"
)

// Options struct for gRPC Server.
type Options struct {
	Address string
	Before  []RequestFunc
	After   []ResponseFunc
	// Other options can be stored in the context.
	Context context.Context
}

// Option function type to handle options.
type Option func(*Options)

// RequestFunc type may take information from the context or request before execute the RPC call.
type RequestFunc func(context.Context, interface{}) (context.Context, interface{})

// ResponseFunc type func may take information from the context or request after execute the RPC call.
type ResponseFunc func(context.Context, interface{}) (context.Context, interface{})

// New returns a Options struct using the opts params.
func New(opts ...Option) *Options {
	options := &Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(options)
	}
	return options
}
