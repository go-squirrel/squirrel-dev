package k8s

import (
	"context"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetNodes() (nodes *v1.NodeList, err error) {
	nodes, err = c.ClientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		zap.S().Error("get node err: ", err)
		return nil, err
	}
	return nodes, err
}

func (c *Client) GetNodeIps() (ips map[string]string, err error) {
	nodes, err := c.GetNodes()
	if err != nil {
		return nil, err
	}
	ips = make(map[string]string, len(nodes.Items))

	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				ips[node.Name] = addr.Address
				break
			}
		}
	}
	return ips, nil
}

func (c *Client) GetNode(nodeName string) (node *v1.Node, err error) {
	node, err = c.ClientSet.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		zap.S().Error("get node err: ", err)
		return nil, err
	}
	return node, err
}

func (c *Client) GetAllNodes() (nodes *v1.NodeList, err error) {
	return c.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) DeleteNode(nodes *v1.NodeList, addresses []string) (err error) {
	nodeNames := []string{}

	for _, address := range addresses {
		for _, node := range nodes.Items {
			for _, addr := range node.Status.Addresses {
				if addr.Type == "InternalIP" {
					if addr.Address == address {
						nodeNames = append(nodeNames, node.Name)
					}
				}
			}
		}
	}

	zap.S().Info("delete nodes: ", nodeNames)
	for _, name := range nodeNames {
		err = c.ClientSet.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			zap.S().Error("delete node:", err)
		}
	}

	return
}

type NodeLabelInfo struct {
	IP     string
	Labels map[string]string
}

func (c *Client) GetNodeLabels() (nodeLabels map[string]NodeLabelInfo, err error) {
	nodes, err := c.GetNodes()
	if err != nil {
		return nil, err
	}

	nodeLabels = make(map[string]NodeLabelInfo, len(nodes.Items))

	for _, node := range nodes.Items {
		nodeIp := ""
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				nodeIp = addr.Address
				break
			}
		}
		nodeLabels[node.Name] = NodeLabelInfo{
			IP:     nodeIp,
			Labels: node.Labels,
		}
	}
	return nodeLabels, nil
}

func (c *Client) SetNodeLabel(node *v1.Node, labels map[string]string) (err error) {
	if node.Labels == nil {
		node.Labels = make(map[string]string)
	}

	for key, value := range labels {
		node.Labels[key] = value
	}
	_, err = c.ClientSet.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	return err
}
