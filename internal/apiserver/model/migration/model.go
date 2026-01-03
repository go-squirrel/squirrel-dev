package migration

import (
	"time"

	"gorm.io/gorm"
)

// Migration 记录数据库迁移版本信息
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"size:50;not null;uniqueIndex"`
	AppliedAt time.Time `gorm:"not null"`
}

// Client 迁移客户端
type Client struct {
	DB *gorm.DB
}

// New 创建迁移客户端
func New(db *gorm.DB) *Client {
	return &Client{
		DB: db,
	}
}

// AutoMigrate 自动迁移版本表结构到数据库
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Migration{})
}

// Init 初始化迁移数据
func Init(db *gorm.DB) {
	// 可以在这里添加初始迁移记录
}

// RecordMigration 记录已应用的迁移版本
func (c *Client) RecordMigration(version string) error {
	has, err := c.HasMigration(version)
	if err != nil {
		return err
	}
	if has {
		return nil
	}
	migration := Migration{
		Version:   version,
		AppliedAt: time.Now(),
	}
	return c.DB.Create(&migration).Error
}

// HasMigration 检查指定版本的迁移是否已应用
func (c *Client) HasMigration(version string) (bool, error) {
	var count int64
	result := c.DB.Model(&Migration{}).Where("version = ?", version).Count(&count)
	return count > 0, result.Error
}

// GetAppliedMigrations 获取所有已应用的迁移版本
func (c *Client) GetAppliedMigrations() ([]Migration, error) {
	var migrations []Migration
	result := c.DB.Order("applied_at asc").Find(&migrations)
	return migrations, result.Error
}
