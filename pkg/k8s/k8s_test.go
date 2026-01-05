package k8s

import (
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateByYaml(t *testing.T) {
	configPath := "./k8s-test-data/kubeconfig"
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err)
	}
	client := New(config)

	yamlContent, err := os.ReadFile("./k8s-test-data/test-deployment.yaml")
	if err != nil {
		panic(err)
	}
	client.CreateByYaml(string(yamlContent))
}
