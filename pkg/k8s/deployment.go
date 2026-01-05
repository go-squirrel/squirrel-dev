package k8s

import (
	"context"

	"go.uber.org/zap"
	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetDeployments(labelSelector string) (deployments *appsV1.DeploymentList, err error) {
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: labelSelector,
		}
	}
	deployments, err = c.ClientSet.AppsV1().Deployments("").List(context.TODO(), listOptions)
	if err != nil {
		zap.S().Error(err)
	}
	return deployments, err
}

func (c *Client) GetDeploymentsByNamespace(namespace string) (deployments *appsV1.DeploymentList, err error) {
	deployments, err = c.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.S().Error(err)
	}
	return deployments, err
}

// 全量的Deployments，根据namespace属性过滤出指定
// 屏蔽多次调用k8sClient
func (c *Client) GetDeploymentsFilterByNamespace(allDeployments *appsV1.DeploymentList, namespace string) (deploymentsRes *appsV1.DeploymentList, err error) {
	// 初始化
	deploymentsRes = &appsV1.DeploymentList{
		TypeMeta: allDeployments.TypeMeta,
		ListMeta: allDeployments.ListMeta,
		Items:    nil,
	}
	for _, deployment := range allDeployments.Items {
		if deployment.Namespace == namespace {
			deploymentsRes.Items = append(deploymentsRes.Items, deployment)
		}
	}
	return
}
