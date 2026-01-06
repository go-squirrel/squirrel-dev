package k8s

// GetKubernetesVersion gets Kubernetes version of the cluster
func (c *Client) GetKubernetesVersion() (string, error) {

	versionInfo, err := c.ClientSet.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	return versionInfo.GitVersion, nil
}
