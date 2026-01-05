package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetNamespaces() (namespaces *v1.NamespaceList, err error) {
	namespaces, err = c.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	return namespaces, err
}

func (c *Client) CreateNamespace(nsName string) (*v1.Namespace, error) {
	// 定义新的命名空间对象
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsName,
		},
	}

	// 创建命名空间
	createdNs, err := c.ClientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace %s: %v", nsName, err)
	}
	fmt.Printf("Successfully created namespace %s\n", nsName)
	return createdNs, nil
}
