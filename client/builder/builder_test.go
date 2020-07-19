package builder_test

import (
	"sync"
	"testing"

	"github.com/aitorfernandez/grpc/client/builder"
	"google.golang.org/grpc/resolver"
)

type testClientConn struct {
	resolver.ClientConn
	target           string
	m1               sync.Mutex
	state            resolver.State
	updateStateCalls int
}

func (t *testClientConn) UpdateState(s resolver.State) {
	t.m1.Lock()
	defer t.m1.Unlock()
	t.state = s
	t.updateStateCalls++
}

func TestBuilder(t *testing.T) {
	target := resolver.Target{
		Scheme:   "test",
		Endpoint: "grpc.io",
	}
	addrs := []string{"localhost:5001", "localhost:5002"}

	b := builder.New(target.Scheme, addrs)
	cc := &testClientConn{target: b.SchemeName}
	r := b.Build(target, cc, resolver.BuildOptions{})

	if r != nil {
		res := r.(*builder.Resolver)

		for i, addr := range res.Addrs[target.Endpoint] {
			if addr != addrs[i] {
				t.Errorf("addrs should be equal")
			}
		}
	}
}
