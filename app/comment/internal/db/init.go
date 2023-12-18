package db

import (
	"context"
	"douyin/config"
	"douyin/model"
	logger2 "douyin/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"strings"
	"time"
)

var DB *gorm.DB

func InitDB() {
	mConfig := config.Conf.MySQL
	host := mConfig.Host
	port := mConfig.Port
	database := mConfig.Database
	username := mConfig.UserName
	password := mConfig.Password
	charset := mConfig.Charset
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=" + charset + "&parseTime=true"}, "")
	err := Database(dsn)
	if err != nil {
		fmt.Println(err)
		logger2.LogrusObj.Error(err)
	}
}

func Database(connString string) error {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connString,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()

	sqlDB.SetConnMaxLifetime(time.Second * 30)
	DB = db
	migration()

	return nil
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := DB
	return db.WithContext(ctx)
}

func migration() {
	// 自动迁移
	err := DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&model.Comment{})
	if err != nil {
		logger2.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	logger2.LogrusObj.Infoln("register table success")
}
