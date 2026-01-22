package script

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"squirrel-dev/internal/pkg/migration"
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

// RegisterMigrations 注册脚本表的迁移
func RegisterMigrations(registry *migration.MigrationRegistry) {
	registry.Register(
		"1.0.0",
		"脚本管理",
		// 升级函数
		func(db *gorm.DB) error {
			err := db.AutoMigrate(&model.Script{})
			if err != nil {
				return err
			}

			// 预置脚本
			scripts := []model.Script{
				{
					Name:    "test-loop",
					Content: "",
				},
			}

			// 读取嵌入的脚本文件并填充内容
			for i := range scripts {
				var content string
				var err error

				switch scripts[i].Name {
				case "test-loop":
					content, err = readFile("test-loop.sh")
				}

				if err != nil {
					return fmt.Errorf("failed to read script for %s: %w", scripts[i].Name, err)
				}
				scripts[i].Content = content
			}

			return db.Create(&scripts).Error
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("scripts")
		},
	)

	registry.Register(
		"1.0.1",
		"脚本执行结果表",
		// 升级函数
		func(db *gorm.DB) error {
			return db.AutoMigrate(&model.ScriptResult{})
		},
		// 回滚函数
		func(db *gorm.DB) error {
			return db.Migrator().DropTable("script_results")
		},
	)
}

// readFile 从嵌入的文件系统中读取脚本文件
func readFile(filename string) (string, error) {
	content, err := fs.ReadFile(ScriptFS, filepath.Join("scripts", filename))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
