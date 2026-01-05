package k8s

import (
	"testing"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestClient_CreatePersistentVolume(t *testing.T) {
	type fields struct {
		ClientSet     *kubernetes.Clientset
		DynamicClient *dynamic.DynamicClient
	}
	type args struct {
		name             string
		size             int64
		storageClassName string
		path             string
		labelKey         string
		labelValue       string
	}
	configPath := "./k8s-test-data/kubeconfig"
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err)
	}
	client := New(config)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				ClientSet: client.ClientSet,
			},
			args: args{
				name:             "test-pv",
				size:             5,
				storageClassName: "test-sc",
				path:             "/test-path",
				labelKey:         "test",
				labelValue:       "0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ClientSet:     tt.fields.ClientSet,
				DynamicClient: tt.fields.DynamicClient,
			}
			if err := c.CreatePersistentVolume(tt.args.name, tt.args.size, tt.args.storageClassName, tt.args.path, tt.args.labelKey, tt.args.labelValue, "ReadWriteMany"); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreatePersistentVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
