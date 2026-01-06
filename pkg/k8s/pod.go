package k8s

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetPodsByNamespace(namespace string) (pods *v1.PodList, err error) {
	if namespace == "all-namespaces" {
		namespace = ""
	}
	pods, err = c.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.S().Error("get pod err: ", err)
		return nil, err
	}
	return pods, err
}

func (c *Client) GetPodsFilterByNamespace(allPods *v1.PodList, namespace string) (podsRes *v1.PodList, err error) {

	podsRes = &v1.PodList{
		TypeMeta: allPods.TypeMeta,
		ListMeta: allPods.ListMeta,
		Items:    nil,
	}
	for _, pod := range allPods.Items {
		if pod.Namespace == namespace {
			podsRes.Items = append(podsRes.Items, pod)
		}
	}
	return
}

func (c *Client) GetAllPods() (pods *v1.PodList, err error) {
	pods, err = c.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.S().Error("get pod err: ", err)
		return nil, err
	}
	return pods, err
}

func isPodReady(pod *v1.Pod) bool {
	for _, containersStatus := range pod.Status.ContainerStatuses {
		fmt.Println("pod:", pod.Name, "/", pod.Namespace, "detail-- container:", containersStatus.Name, "status:", containersStatus.Ready)
		if containersStatus.State.Terminated != nil && containersStatus.State.Terminated.Reason == "Completed" {
			zap.S().Warnf("Pod %s/%s is Completed, pass", pod.Namespace, pod.Name)
			return true
		} else if !containersStatus.Ready {
			zap.S().Warnf("Pod %s/%s is not ready", pod.Namespace, pod.Name)
			return false
		}
	}
	zap.S().Infof("Pod %s/%s is ready", pod.Namespace, pod.Name)
	return true
}

func (c *Client) CheckPodsReadiness(namespace string) error {
	deployments, err := c.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to get Deployments: %v", err)
	}

	for _, deployment := range deployments.Items {
		selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
		if err != nil {
			return fmt.Errorf("Failed to convert LabelSelector: %v", err)
		}
		pods, err := c.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			return fmt.Errorf("Failed to get Pods: %v", err)
		}
		for _, pod := range pods.Items {
			if !isPodReady(&pod) {
				zap.S().Warnf("Pod %s/%s associated with Deployment is not ready", pod.Namespace, pod.Name)
			}
		}
	}

	statefulSets, err := c.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to get StatefulSets: %v", err)
	}

	// For each StatefulSet, get its Pods
	for _, sts := range statefulSets.Items {
		selector, err := metav1.LabelSelectorAsSelector(sts.Spec.Selector)
		if err != nil {
			return fmt.Errorf("Failed to convert LabelSelector: %v", err)
		}
		pods, err := c.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			return fmt.Errorf("Failed to get Pods: %v", err)
		}
		for _, pod := range pods.Items {
			if !isPodReady(&pod) {
				zap.S().Warnf("Pod %s/%s associated with StatefulSet is not ready", pod.Namespace, pod.Name)
			}
		}
	}
	return nil
}

func (c *Client) DeletePod(name, namespace string) (err error) {
	err = c.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}
