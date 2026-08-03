package application

import (
	"context"
	"errors"
	"testing"

	"squirrel-dev/internal/squ-agent/module/server/domain"
)

type fakeCollector struct {
	value *domain.HostInfo
	err   error
}

func (f fakeCollector) CollectHostInfo(context.Context) (*domain.HostInfo, error) {
	return f.value, f.err
}

func TestGetInfo(t *testing.T) {
	expected := &domain.HostInfo{
		Hostname: "legacy-host",
		OS:       "linux",
		IPAddresses: []domain.NetAddr{{
			Name: "eth0",
			IPv4: []string{"192.0.2.10"},
		}},
	}
	service := NewService(fakeCollector{value: expected})

	actual, err := service.GetInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("result = %#v, want %#v", actual, expected)
	}
}

func TestGetInfoFailure(t *testing.T) {
	expected := errors.New("collect failed")
	service := NewService(fakeCollector{err: expected})

	_, err := service.GetInfo(context.Background())
	if !errors.Is(err, expected) {
		t.Fatalf("error = %v, want %v", err, expected)
	}
}

func TestGetInfoWithoutCollector(t *testing.T) {
	_, err := NewService(nil).GetInfo(context.Background())
	if !errors.Is(err, ErrCollectorUnavailable) {
		t.Fatalf("error = %v, want %v", err, ErrCollectorUnavailable)
	}
}
