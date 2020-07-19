package grpc_test

import (
	"context"
	"fmt"
	"log"

	"testing"
	"time"

	b "github.com/aitorfernandez/grpc/client/balancer"
	"github.com/aitorfernandez/grpc/client/builder"
	pb "github.com/aitorfernandez/grpc/proto"
	"github.com/aitorfernandez/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

const (
	scheme      = "foo"
	serviceName = "customBalancer"
)

var (
	addrs           = []string{":50051", ":50052", ":50053"}
	balancingPolicy = b.Name
)

func init() {
	resolver.Register(builder.New(scheme, addrs))
	balancer.Register(b.NewBuilder())
}

type testServer struct {
	pb.UnimplementedTestServer
	Addr string
}

func (s *testServer) Send(ctx context.Context, req *pb.Req) (*pb.Res, error) {
	return &pb.Res{
		Pong: fmt.Sprintf("%s and Pong from %s", req.GetPing(), s.Addr),
	}, nil
}

func callSend(c pb.TestClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, _ := c.Send(ctx, &pb.Req{Ping: "Ping"})

	return res.GetPong()
}

func makeCalls(c *grpc.ClientConn, n int) []string {
	var res []string

	tc := pb.NewTestClient(c)
	for i := 0; i < n; i++ {
		res = append(res, callSend(tc))
	}

	return res
}

func startServe(addr string) server.Server {
	s := server.New(
		server.Address(addr),
	)

	pb.RegisterTestServer(s.GRPC(), &testServer{
		Addr: s.Options().Address,
	})

	if err := s.Start(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}

	return s
}

func TestGRPC(t *testing.T) {
	var ss []server.Server
	for _, addr := range addrs {
		ss = append(ss, startServe(addr))
	}

	for _, s := range ss {
		defer func(s server.Server) {
			s.Stop()
		}(s)
	}

	nodeBalancerConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", scheme, serviceName),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, balancingPolicy)),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err := nodeBalancerConn.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	res := makeCalls(nodeBalancerConn, len(addrs))

	for _, addr := range addrs {
		msg := fmt.Sprintf("Ping and Pong from %s", addr)
		if !in(res, msg) {
			t.Errorf("msg should be in the response, %w %s", res, msg)
		}

		if total := count(res, msg); total != 1 {
			t.Errorf("msg should be unique in the response")
		}
	}

	if len(res) != len(addrs) {
		t.Error("different number of client responses")
	}
}

func in(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func count(slice []string, val string) int {
	var i int
	for _, item := range slice {
		if item == val {
			i++
		}
	}
	return i
}
