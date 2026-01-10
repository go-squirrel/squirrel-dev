package collector_test

import (
	"encoding/json"
	"fmt"
	"squirrel-dev/pkg/collector"
	"testing"
)

func TestCollectorFactory_CollectAll(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    *collector.HostInfo
		wantErr bool
	}{
		{name: "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := collector.NewCollectorFactory()
			f.Register(collector.NewCPUCollector())
			f.Register(collector.NewMemoryCollector())
			f.Register(collector.NewDiskCollector())
			got, _ := f.CollectAll()
			data, _ := json.Marshal(got)
			fmt.Println(string(data))
		})
	}
}
