package app_store

import (
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

type AppStore struct {
	gorm.Model
	Name        string  `gorm:"column:name;type:varchar(100);not null;unique;comment:模板名称"`
	Description string  `gorm:"column:description;type:text;comment:模板描述"`
	Type        string  `gorm:"column:type;type:varchar(50);not null;comment:类型(compose/k8s_manifest/helm_chart)"`
	Category    string  `gorm:"column:category;type:varchar(50);comment:分类(web/database/middleware/devops)"`
	Icon        *string `gorm:"column:icon;type:varchar(255);comment:图标URL"`
	Version     string  `gorm:"column:version;type:varchar(50);not null;comment:版本"`
	Content     string  `gorm:"column:content;type:text;not null;comment:模板内容"`
	Tags        string  `gorm:"column:tags;type:varchar(255);comment:标签(逗号分隔)"`
	Author      string  `gorm:"column:author;type:varchar(100);comment:作者"`
	RepoUrl     *string `gorm:"column:repo_url;type:varchar(255);comment:仓库地址"`
	HomepageUrl *string `gorm:"column:homepage_url;type:varchar(255);comment:主页地址"`
	IsOfficial  bool    `gorm:"column:is_official;type:tinyint(1);default:false;comment:是否官方模板"`
	Downloads   int     `gorm:"column:downloads;type:int(10);default:0;comment:下载次数"`
	Status      string  `gorm:"column:status;type:varchar(20);default:'active';comment:状态(active/deprecated)"`
}

func (AppStore) TableName() string {
	return "app_stores"
}
