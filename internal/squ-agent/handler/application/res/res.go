package res

type Application struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Content     string `json:"content"`
	Version     string `json:"version"`
}

