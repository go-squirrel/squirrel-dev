package queue

import (
	"squirrel-dev/internal/pkg/database"

	"go.uber.org/zap"
)

type Queue struct {
	DB database.DB
}

func New(db database.DB) *Queue {
	return &Queue{
		DB: db,
	}
}

// Start 启动队列处理器
func (q *Queue) Start() {
	zap.L().Info("队列处理器已启动")
	// 队列处理逻辑在其他地方实现，这里只是占位
}
