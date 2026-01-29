package utils_test

import (
	"fmt"
	"squirrel-dev/pkg/utils"

	"testing"
)

func TestGenAgentUrl(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		schema  string
		address string
		port    int
		baseuri string
		uri     string
		want    string
	}{
		{
			name:    "http",
			schema:  "http",
			address: "127.0.0.1",
			port:    0,
			baseuri: "/api/v1",
			uri:     "application",
			want:    "http://127.0.0.1/api/v1/application",
		},
		{
			name:    "https",
			schema:  "https",
			address: "127.0.0.1",
			port:    10443,
			baseuri: "/api/v1",
			uri:     "application",
			want:    "https://127.0.0.1:10443/api/v1/application",
		},
		{
			name:    "https",
			schema:  "https",
			address: "example.com:10444",
			port:    10443,
			baseuri: "/api/v1",
			uri:     "application",
			want:    "https://example.com:10444/api/v1/application",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.GenAgentUrl(tt.schema, tt.address, tt.port, tt.baseuri, tt.uri)
			if got == tt.want {
				fmt.Println("got:", got)
			} else {
				t.Errorf("GenAgentUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
