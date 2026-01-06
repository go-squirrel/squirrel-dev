package k8s

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Secret struct definition
type Secret struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   SecretMetadata    `yaml:"metadata"`
	Type       string            `yaml:"type"`
	Data       map[string]string `yaml:"data"`
}

// SecretMetadata contains Secret metadata
type SecretMetadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}

// Function to create TLS certificate Secret
func (c *Client) GenerateTLSSecret(
	name string,
	namespace string,
	tlsCrt string, // Base64 encoded certificate
	tlsKey string, // Base64 encoded private key
	labelKey string,
	labelValue string,
	caCrt ...string,
) (secret *Secret, err error) {

	data := map[string]string{
		"tls.crt": tlsCrt,
		"tls.key": tlsKey,
	}
	if len(caCrt) > 0 {
		data["ca.crt"] = caCrt[0]
	}

	secret = &Secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Metadata: SecretMetadata{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				labelKey: labelValue,
			},
		},
		Type: "kubernetes.io/tls", // TLS certificate type
		Data: data,
	}
	return secret, nil
}

// dataFile: key-filePath
func (c *Client) CreateSecretByFile(namespace, name string, dataFile map[string]string, secretType v1.SecretType) (err error) {

	data := map[string]string{}
	for key, value := range dataFile {
		file, err := os.ReadFile(value)
		if err != nil {
			fmt.Println("Err open ", value)
		}
		data[key] = string(file)
	}
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: secretType,
		Data: make(map[string][]byte),
	}

	for k, v := range data {
		secret.Data[k] = []byte(v)
	}

	_, err = c.ClientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	return err
}
