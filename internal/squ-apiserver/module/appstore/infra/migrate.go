package infra

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"

	"gorm.io/gorm"
)

//go:embed templates/*.yml
var templates embed.FS

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&appModel{}); err != nil {
		return err
	}
	pointer := func(value string) *string { return &value }
	apps := []appModel{
		{Name: "Nginx", Description: "web server.", Type: "compose", Category: "web", Version: "1.28.0", Author: "Elastic", HomepageURL: pointer("https://nginx.org/"), RepoURL: pointer("https://github.com/nginx/nginx"), IsOfficial: true, Status: "active"},
		{Name: "MySQL", Description: "MySQL 8.0 with persistence, monitoring and exporter. Open-source relational database.", Type: "compose", Category: "database", Version: "8.0.0", Author: "Oracle", HomepageURL: pointer("https://www.mysql.com/"), RepoURL: pointer("https://github.com/mysql/mysql-server"), IsOfficial: true, Status: "active"},
		{Name: "Redis", Description: "Redis 7.2 with persistence, Redis Commander and monitoring. In-memory data structure store.", Type: "compose", Category: "database", Version: "7.2.0", Author: "Redis Labs", HomepageURL: pointer("https://redis.io/"), RepoURL: pointer("https://github.com/redis/redis"), IsOfficial: true, Status: "active"},
		{Name: "Elasticsearch", Description: "Elasticsearch 8.11 with Kibana and exporter. Distributed search and analytics engine.", Type: "compose", Category: "database", Version: "8.11.0", Author: "Elastic", HomepageURL: pointer("https://www.elastic.co/"), RepoURL: pointer("https://github.com/elastic/elasticsearch"), IsOfficial: true, Status: "active"},
		{Name: "Jenkins", Description: "Jenkins 2.401.1 with persistence, monitoring and exporter. Continuous integration and delivery tool.", Type: "compose", Category: "database", Version: "latest", Author: "Jenkins", HomepageURL: pointer("https://www.jenkins.io/"), RepoURL: pointer("https://github.com/jenkinsci/jenkins"), IsOfficial: true, Status: "active"},
	}
	for i := range apps {
		filename := map[string]string{"Nginx": "nginx-compose.yml", "MySQL": "mysql-compose.yml", "Redis": "redis-compose.yml", "Elasticsearch": "elasticsearch-compose.yml", "Jenkins": "jenkins-compose.yml"}[apps[i].Name]
		content, err := fs.ReadFile(templates, filepath.Join("templates", filename))
		if err != nil {
			return fmt.Errorf("failed to read template for %s: %w", apps[i].Name, err)
		}
		apps[i].Content = string(content)
	}
	return db.Create(&apps).Error
}

func Rollback(db *gorm.DB) error { return db.Migrator().DropTable("app_stores") }
