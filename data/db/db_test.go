package db

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	database := GetDb()

	// 尝试 ping 数据库
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("Failed to get sql.DB from gorm.DB: %s", err.Error())
	}

	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %s", err.Error())
	}

	t.Log("Database connection successful")
}
