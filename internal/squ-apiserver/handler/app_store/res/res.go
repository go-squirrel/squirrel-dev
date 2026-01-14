package res

type AppStore struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Icon        *string `json:"icon"`
	Version     string  `json:"version"`
	Content     string  `json:"content"`
	Tags        string  `json:"tags"`
	Author      string  `json:"author"`
	RepoUrl     *string `json:"repo_url"`
	HomepageUrl *string `json:"homepage_url"`
	IsOfficial  bool    `json:"is_official"`
	Downloads   int     `json:"downloads"`
	Status      string  `json:"status"`
}
