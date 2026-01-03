package database

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SQLiteDB struct {
	conn *gorm.DB
}

func (db *SQLiteDB) Connect(connectionString string) error {
	var err error
	dir := filepath.Dir(connectionString)

	// 如果目录非空（不是当前目录 "."），则创建它
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	db.conn, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	return err
}

func (db *SQLiteDB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *SQLiteDB) GetDB() *gorm.DB {
	return db.conn
}
