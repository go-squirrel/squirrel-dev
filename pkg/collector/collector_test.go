package collector

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Collect(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    any
		wantErr bool
	}{
		{
			name: "memory",
		},
		{
			name: "disk",
		},
		{
			name: "cpu",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				ok   bool
				data []byte
				info any
				got  any
			)
			switch tt.name {
			case "memory":
				c := NewMemoryCollector()
				got, _ := c.Collect()
				data, _ = json.Marshal(got)
				info, ok = got.(*MemInfo)
			case "disk":
				c := NewDiskCollector()
				got, _ := c.Collect()
				data, _ = json.Marshal(got)
				info, ok = got.(*DiskInfo)
			case "cpu":
				c := NewCPUCollector()
				got, _ := c.Collect()
				data, _ = json.Marshal(got)
				info, ok = got.(*CPUInfo)
			}
			if !ok {
				t.Errorf("Collect() = %v, want %v", got, tt.want)
			}
			fmt.Println(string(data))
			if !ok {
				t.Errorf("Collect() = %v, want %v", got, tt.want)
			}
			fmt.Println(info)

		})
	}
}
