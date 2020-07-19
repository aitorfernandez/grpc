package option_test

import (
	"testing"

	"github.com/aitorfernandez/grpc/server/option"
)

func Address(a string) option.Option {
	return func(o *option.Options) {
		o.Address = a
	}
}

func TestNew(t *testing.T) {
	addr := ":50051"
	got := option.New(Address(addr))

	if got.Address != addr {
		t.Errorf("got %v want %s", got, addr)
	}

	if got.Context == nil {
		t.Errorf("got %v want context", got.Context)
	}
}
