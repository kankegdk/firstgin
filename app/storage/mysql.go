package storage

import (
	"database/sql"
	"fmt"
	"log"
	"myapi/app/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "github.com/go-sql-driver/mysql"
)

var MySQLDB *sql.DB
var GormDB *gorm.DB

// InitMySQL 初始化MySQL连接池
func InitMySQL() error {
	cfg := config.GetDatabaseConfig()

	// 构建DSN (Data Source Name)
	dsn := cfg.GetConnectionString()

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime))

	MySQLDB = db
	log.Println("MySQL连接池初始化成功")

	// 初始化GORM
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize GORM: %v", err)
	}

	GormDB = gormDB
	log.Println("GORM初始化成功")
	return nil
}

// GetMySQL 获取MySQL连接实例
func GetMySQL() *sql.DB {
	return MySQLDB
}

// GetGormDB 获取GORM连接实例
func GetGormDB() *gorm.DB {
	return GormDB
}

// CloseMySQL 关闭MySQL连接
func CloseMySQL() {
	if MySQLDB != nil {
		MySQLDB.Close()
		log.Println("MySQL连接已关闭")
	}
}
