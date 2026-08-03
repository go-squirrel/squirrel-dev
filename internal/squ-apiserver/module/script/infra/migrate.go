package infra

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"gorm.io/gorm"
)

func MigrateScripts(db *gorm.DB) error {
	if err := db.AutoMigrate(&scriptModel{}); err != nil {
		return err
	}
	content, err := readScript("test-loop.sh")
	if err != nil {
		return fmt.Errorf("failed to read script for test-loop: %w", err)
	}
	return db.Create(&[]scriptModel{{Name: "test-loop", Content: content}}).Error
}

func RollbackScripts(db *gorm.DB) error {
	return db.Migrator().DropTable("scripts")
}

func MigrateResults(db *gorm.DB) error {
	return db.AutoMigrate(&resultModel{})
}

func RollbackResults(db *gorm.DB) error {
	return db.Migrator().DropTable("script_results")
}

func readScript(filename string) (string, error) {
	content, err := fs.ReadFile(scriptFS, filepath.Join("scripts", filename))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
