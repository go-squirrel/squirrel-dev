package k8s

import (
	"context"
	"strings"

	"go.uber.org/zap"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetStatefulsets() (statefulsets *appsV1.StatefulSetList, err error) {
	statefulsets, err = c.ClientSet.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.S().Error(err)
	}
	return statefulsets, err
}

func (c *Client) GetStatefulsetsByNamespace(namespace string) (statefulsets *appsV1.StatefulSetList, err error) {
	statefulsets, err = c.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.S().Error(err)
	}
	return statefulsets, err
}

func (c *Client) GetPodByStatefulsetName(pods *v1.PodList, statefulSetName string) (statefulsetPods []*v1.Pod, err error) {
	for _, pod := range pods.Items {
		// 检查 Pod 的标签，看是否属于指定的 StatefulSet
		if val, ok := pod.Labels["statefulset.kubernetes.io/pod-name"]; ok && strings.HasPrefix(val, statefulSetName) {
			podCopy := pod
			statefulsetPods = append(statefulsetPods, &podCopy)
		}
	}
	return statefulsetPods, err
}
