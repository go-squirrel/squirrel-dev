package monitor

import (
	"gorm.io/gorm"
)

type MonitorClient struct {
	DB *gorm.DB
}
