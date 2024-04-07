// Package main TODO
package main

import (
	"context"
	"os"

	"github.com/onesaltedseafish/go-utils/log"
	gormlog "github.com/onesaltedseafish/go-utils/log/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ctx    = context.Background()
	logOpt = log.CommonLogOpt.WithDirectory("/tmp/gorm").WithLogLevel(zapcore.DebugLevel)
)

// User 结构体
type User struct {
	ID   uint
	Name string
	Age  uint
}

func main() {
	dbLogger := gormlog.NewLogger("gorm", &logOpt)
	mainLogger := log.GetLogger("main", &logOpt)

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		mainLogger.Error(ctx, "Failed to connect to database", zap.Error(err))
		os.Exit(-1)
	}

	// 自动迁移 User 结构体对应的表
	err = db.AutoMigrate(&User{})
	if err != nil {
		mainLogger.Error(ctx, "Failed to auto migrate", zap.Error(err))
		os.Exit(-1)
	}

	// 创建新的用户
	user := User{Name: "Alice", Age: 30}
	result := db.Create(&user)
	if result.Error != nil {
		mainLogger.Error(ctx, "Failed to create user", zap.Error(result.Error))
		os.Exit(-1)
	}

	mainLogger.Info(ctx, "Created user", zap.Any("user", user))

	// 查询用户
	var fetchedUser User
	db.First(&fetchedUser, user.ID)
	mainLogger.Info(ctx, "Fetched user", zap.Any("user", fetchedUser))
}
