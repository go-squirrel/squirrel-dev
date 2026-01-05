package k8s

// getKubernetesVersion 获取集群的 Kubernetes 版本
func (c *Client) GetKubernetesVersion() (string, error) {

	versionInfo, err := c.ClientSet.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	return versionInfo.GitVersion, nil
}
