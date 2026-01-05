package k8s

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 根据serviceName获取clusterIp
func (c *Client) GetClusterIPByName(serviceName string, namespace string) (clusterIp string, err error) {
	service, err := c.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), serviceName, v1.GetOptions{})
	if err != nil {
		return "", err
	}
	clusterIp = service.Spec.ClusterIP
	return clusterIp, nil
}
