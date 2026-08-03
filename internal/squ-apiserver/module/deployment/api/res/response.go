package res

// ApplicationInfo describes the deployed application.
type ApplicationInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Version     string `json:"version"`
}

// ServerInfo describes the deployment target.
type ServerInfo struct {
	ID        uint   `json:"id"`
	IPAddress string `json:"ip_address"`
	AgentPort int    `json:"agent_port"`
}

// Deployment describes an application deployment.
type Deployment struct {
	ID          uint            `json:"id"`
	DeployID    uint64          `json:"deploy_id"`
	Application ApplicationInfo `json:"application"`
	Server      ServerInfo      `json:"server"`
	Status      string          `json:"status"`
	DeployedAt  string          `json:"deployed_at"`
	Content     string          `json:"content"`
}
