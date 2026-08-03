package req

// DeployApplication is the request to deploy an application to a server.
type DeployApplication struct {
	ServerID uint `json:"server_id"`
}

// UpdateDeployment is the request to update deployment content.
type UpdateDeployment struct {
	Content string `json:"content"`
}

// ReportApplicationStatus is the request used by an agent to report deployment status.
type ReportApplicationStatus struct {
	ApplicationID uint   `json:"application_id"`
	ServerID      uint   `json:"server_id"`
	Status        string `json:"status"`
	DeployID      uint64 `json:"deploy_id"`
}
