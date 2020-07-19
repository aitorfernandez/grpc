package balancer

import (
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// Name is the name of the custom balancer.
const Name = "my_custom_balancer"

// sentConn keeps track of the connection sent.
var sentConn []balancer.SubConn

// NewBuilder creates a new huski balancer builder.
func NewBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(Name, &nodePickerBuilder{}, base.Config{HealthCheck: true})
}

type nodePickerBuilder struct{}

// Build returns a nodePicker each time the builder state is updated.
func (*nodePickerBuilder) Build(info base.PickerBuildInfo) balancer.V2Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	var scs []balancer.SubConn
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}

	return &nodePicker{
		subConns: scs,
	}
}

type nodePicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
}

// Pick keeps the connections and return a connection if it has not been used before.
func (p *nodePicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()

	if len(sentConn) == len(p.subConns) {
		sentConn = sentConn[:0]
	}

	var sc balancer.SubConn
	for _, s := range p.subConns {
		if in(sentConn, s) {
			continue
		}
		sc = s
	}

	sentConn = append(sentConn, sc)

	p.mu.Unlock()

	return balancer.PickResult{SubConn: sc}, nil
}

func in(scs []balancer.SubConn, n balancer.SubConn) bool {
	for _, sc := range scs {
		if sc == n {
			return true
		}
	}
	return false
}
