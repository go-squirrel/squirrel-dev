package k8s

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestClient_CreateSecretByFile(t *testing.T) {
	type fields struct {
		ClientSet     *kubernetes.Clientset
		DynamicClient *dynamic.DynamicClient
	}
	type args struct {
		namespace  string
		name       string
		dataFile   map[string]string
		secretType v1.SecretType
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
				namespace: "sd-cloud",
				name:      "secret-client",
				dataFile: map[string]string{
					"client.crt": "./k8s-test-data/certs/client.crt",
					"client.key": "./k8s-test-data/certs/client.key",
					"ca.crt":     "./k8s-test-data/certs/ca.crt",
				},
				secretType: v1.SecretTypeOpaque,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ClientSet:     tt.fields.ClientSet,
				DynamicClient: tt.fields.DynamicClient,
			}
			c.CreateSecretByFile(tt.args.namespace, tt.args.name, tt.args.dataFile, tt.args.secretType)
		})
	}
}

func TestClient_GenerateTLSSecret(t *testing.T) {
	type fields struct {
		ClientSet     *kubernetes.Clientset
		DynamicClient *dynamic.DynamicClient
	}
	type args struct {
		name       string
		namespace  string
		tlsCrt     string
		tlsKey     string
		labelKey   string
		labelValue string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantSecret *Secret
		wantErr    bool
	}{
		{
			args: args{
				name:       "secret-client",
				namespace:  "sd-cloud",
				tlsCrt:     "./k8s-test-data/certs/client.crt",
				tlsKey:     "./k8s-test-data/certs/client.key",
				labelKey:   "app",
				labelValue: "client",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ClientSet:     tt.fields.ClientSet,
				DynamicClient: tt.fields.DynamicClient,
			}
			gotSecret, err := c.GenerateTLSSecret(tt.args.name, tt.args.namespace, tt.args.tlsCrt, tt.args.tlsKey, tt.args.labelKey, tt.args.labelValue)
			if err != nil {
				t.Errorf("GenerateTLSSecret() error = %v", err)
				return
			}
			yamlString, err := yaml.Marshal(gotSecret)
			if err != nil {
				t.Errorf("GenerateTLSSecret() error = %v", err)
				return
			}
			fmt.Println(string(yamlString))
		})
	}
}
