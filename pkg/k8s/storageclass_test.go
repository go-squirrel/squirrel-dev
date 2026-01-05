package k8s

import (
	"testing"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestClient_CreateStorageClass(t *testing.T) {
	type fields struct {
		ClientSet     *kubernetes.Clientset
		DynamicClient *dynamic.DynamicClient
	}
	type args struct {
		name string
	}

	configPath := "./k8s-test-data/kubeconfig"
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err)
	}
	client := New(config)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				ClientSet: client.ClientSet,
			},
			args: args{
				name: "local-storage",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ClientSet:     tt.fields.ClientSet,
				DynamicClient: tt.fields.DynamicClient,
			}
			c.CreateStorageClass(tt.args.name)
		})
	}
}
