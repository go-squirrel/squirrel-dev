package migration

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// MigrationFunc 定义迁移函数类型
type MigrationFunc func(*gorm.DB) error

// MigrationItem 表示一个迁移项
type MigrationItem struct {
	Version string
	Name    string
	Up      MigrationFunc
	Down    MigrationFunc
}

// MigrationRegistry 迁移注册表
type MigrationRegistry struct {
	Migrations []MigrationItem
}

// NewMigrationRegistry 创建新的迁移注册表
func NewMigrationRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		Migrations: make([]MigrationItem, 0),
	}
}

// Register 注册一个迁移
func (r *MigrationRegistry) Register(version, name string, up, down MigrationFunc) {
	r.Migrations = append(r.Migrations, MigrationItem{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

// RunMigrations 运行所有未应用的迁移
func RunMigrations(db *gorm.DB, registry *MigrationRegistry) error {
	// 确保迁移表存在
	AutoMigrate(db)

	// 创建迁移客户端
	client := New(db)

	// 获取已应用的迁移
	appliedMigrations, err := client.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("获取已应用迁移失败: %w", err)
	}

	// 创建已应用迁移的映射，方便查找
	appliedMap := make(map[string]bool)
	for _, m := range appliedMigrations {
		appliedMap[m.Version] = true
	}

	// 运行未应用的迁移
	for _, migration := range registry.Migrations {
		if !appliedMap[migration.Version] {
			log.Printf("应用迁移: %s - %s", migration.Version, migration.Name)

			// 先在事务内执行迁移
			tx := db.Begin()
			if tx.Error != nil {
				return fmt.Errorf("开始事务失败: %w", tx.Error)
			}

			if err := migration.Up(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("迁移失败 %s: %w", migration.Version, err)
			}

			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("提交事务失败: %w", err)
			}

			// 事务外记录迁移，避免锁冲突
			if err := client.RecordMigration(migration.Version); err != nil {
				return fmt.Errorf("记录迁移失败 %s: %w", migration.Version, err)
			}

			log.Printf("迁移成功: %s", migration.Version)
		}
	}

	return nil
}

// RollbackMigration 回滚指定版本的迁移
func RollbackMigration(db *gorm.DB, registry *MigrationRegistry, version string) error {
	// 确保迁移表存在
	AutoMigrate(db)

	// 创建迁移客户端
	client := New(db)

	// 检查迁移是否已应用
	hasMigration, err := client.HasMigration(version)
	if err != nil {
		return fmt.Errorf("检查迁移状态失败: %w", err)
	}

	if !hasMigration {
		return fmt.Errorf("迁移 %s 未应用，无法回滚", version)
	}

	// 查找迁移项
	var migrationItem *MigrationItem
	for _, m := range registry.Migrations {
		if m.Version == version {
			migrationItem = &m
			break
		}
	}

	if migrationItem == nil {
		return fmt.Errorf("未找到迁移 %s 的定义", version)
	}

	log.Printf("回滚迁移: %s - %s", migrationItem.Version, migrationItem.Name)

	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %w", tx.Error)
	}

	// 执行回滚
	if err := migrationItem.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("回滚失败 %s: %w", version, err)
	}

	// 删除迁移记录
	if err := tx.Where("version = ?", version).Delete(&Migration{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除迁移记录失败 %s: %w", version, err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	log.Printf("回滚成功: %s", version)
	return nil
}
