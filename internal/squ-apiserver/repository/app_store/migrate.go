package app_store

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

func RegisterMigrations(registry *migration.MigrationRegistry) {
	// 注册应用商店表迁移
	registry.Register(
		"1.0.0",
		"应用商店",
		// 升级函数
		func(db *gorm.DB) error {
			err := db.AutoMigrate(&model.AppStore{})
			if err != nil {
				return err
			}

			// 预置应用模板
			apps := []model.AppStore{
				{
					Name:        "Nginx",
					Description: "web server.",
					Type:        TypeCompose,
					Category:    CategoryWeb,
					Version:     "1.28.0",
					Author:      "Elastic",
					HomepageUrl: strPtr("https://nginx.org/"),
					RepoUrl:     strPtr("https://github.com/nginx/nginx"),
					IsOfficial:  true,
					Downloads:   0,
					Status:      StatusActive,
				},
				{
					Name:        "MySQL",
					Description: "MySQL 8.0 with persistence, monitoring and exporter. Open-source relational database.",
					Type:        TypeCompose,
					Category:    CategoryDatabase,
					Version:     "8.0.0",
					Author:      "Oracle",
					HomepageUrl: strPtr("https://www.mysql.com/"),
					RepoUrl:     strPtr("https://github.com/mysql/mysql-server"),
					IsOfficial:  true,
					Downloads:   0,
					Status:      StatusActive,
				},
				{
					Name:        "Redis",
					Description: "Redis 7.2 with persistence, Redis Commander and monitoring. In-memory data structure store.",
					Type:        TypeCompose,
					Category:    CategoryDatabase,
					Version:     "7.2.0",
					Author:      "Redis Labs",
					HomepageUrl: strPtr("https://redis.io/"),
					RepoUrl:     strPtr("https://github.com/redis/redis"),
					IsOfficial:  true,
					Downloads:   0,
					Status:      StatusActive,
				},
				{
					Name:        "Elasticsearch",
					Description: "Elasticsearch 8.11 with Kibana and exporter. Distributed search and analytics engine.",
					Type:        TypeCompose,
					Category:    CategoryDatabase,
					Version:     "8.11.0",
					Author:      "Elastic",
					HomepageUrl: strPtr("https://www.elastic.co/"),
					RepoUrl:     strPtr("https://github.com/elastic/elasticsearch"),
					IsOfficial:  true,
					Downloads:   0,
					Status:      StatusActive,
				},
				{
					Name:        "Jenkins",
					Description: "Jenkins 2.401.1 with persistence, monitoring and exporter. Continuous integration and delivery tool.",
					Type:        TypeCompose,
					Category:    CategoryDatabase,
					Version:     "latest",
					Author:      "Jenkins",
					HomepageUrl: strPtr("https://www.jenkins.io/"),
					RepoUrl:     strPtr("https://github.com/jenkinsci/jenkins"),
					IsOfficial:  true,
					Downloads:   0,
					Status:      StatusActive,
				},
			}

			// 读取嵌入的模板文件并填充内容
			for i := range apps {
				var content string
				var err error

				switch apps[i].Name {
				case "MySQL":
					content, err = readFile("mysql-compose.yml")
				case "Redis":
					content, err = readFile("redis-compose.yml")
				case "Elasticsearch":
					content, err = readFile("elasticsearch-compose.yml")
				case "Nginx":
					content, err = readFile("nginx-compose.yml")
				case "Jenkins":
					content, err = readFile("jenkins-compose.yml")
				}

				if err != nil {
					return fmt.Errorf("failed to read template for %s: %w", apps[i].Name, err)
				}
				apps[i].Content = content
			}

			return db.Create(&apps).Error
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("app_stores")
		},
	)
}

// readFile 从嵌入的文件系统中读取文件内容
func readFile(filename string) (string, error) {
	content, err := fs.ReadFile(TemplateFS, filepath.Join("templates", filename))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// strPtr 返回字符串指针的辅助函数
func strPtr(s string) *string {
	return &s
}
