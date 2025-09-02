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
	// 直接获取数据库配置
	user := config.GetString("dbUser", "root")
	password := config.GetString("dbPassword", "")
	host := config.GetString("dbHost", "localhost")
	port := config.GetString("dbPort", "3306")
	dbName := config.GetString("dbName", "mydatabase")
	charset := config.GetString("dbCharset", "utf8mb4")
	maxLifetime := config.GetInt("dbMaxLifetime", 30)
	maxIdleConns := config.GetInt("dbMaxIdleConns", 100)
	maxOpenConns := config.GetInt("dbMaxOpenConns", 10)

	// 构建DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		user, password, host, port, dbName, charset)

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
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifetime))

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
