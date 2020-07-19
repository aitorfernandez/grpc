package builder

import (
	"google.golang.org/grpc/resolver"
)

// New creates a ResolverBuilder.
func New(scheme string, addrs []string) *ResolverBuilder {
	return &ResolverBuilder{
		SchemeName: scheme,
		Addrs:      addrs,
	}
}

// ResolverBuilder will be registered to watch the updates for the target and send updates to the ClientConn.
type ResolverBuilder struct {
	SchemeName string
	Addrs      []string
}

// Build creates and returns a new resolver for the given target.
func (b ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		Target: target,
		CConn:  cc,
		Addrs: map[string][]string{
			target.Endpoint: b.Addrs,
		},
	}

	r.Update()
	return r, nil
}

// Scheme returns the scheme supported.
func (b ResolverBuilder) Scheme() string {
	return b.SchemeName
}

// Resolver struct for address updates. https://godoc.org/google.golang.org/grpc/resolver#Resolver
type Resolver struct {
	Target resolver.Target
	CConn  resolver.ClientConn
	Addrs  map[string][]string
}

// Update updates the addrs.
func (r Resolver) Update() {
	aa := r.Addrs[r.Target.Endpoint]
	addrs := make([]resolver.Address, len(aa))
	for i, a := range aa {
		addrs[i] = resolver.Address{Addr: a}
	}

	// Manually provide set of resolved addresses for the target.
	r.CConn.UpdateState(resolver.State{Addresses: addrs})
}

// ResolveNow https://godoc.org/google.golang.org/grpc/resolver#Resolver
func (Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close https://godoc.org/google.golang.org/grpc/resolver#Resolver
func (Resolver) Close() {}
