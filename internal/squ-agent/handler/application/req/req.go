package req

type Application struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
	Status      string              `json:"status"`
	Content     string              `json:"content"`
	Version     string              `json:"version"`
	ServerID    uint                `json:"server_id"`
	DeployID    uint64              `json:"deploy_id"`
	Env         []map[string]string `json:"env"`
}
