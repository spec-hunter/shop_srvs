package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ahlemarg/shop-srvs/src/user_srvs/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	mysqlInfo := global.ServerInfo.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlInfo.User, mysqlInfo.PassWord, mysqlInfo.Host, mysqlInfo.Port, mysqlInfo.DB)

	logFilePath := "./logs/db_logs/db.logs"
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		res := fmt.Sprintf("打开文件出现错误: %v.", err)
		fmt.Println(res)
	}

	NewLogger := logger.New(
		log.New(file, "\r\t", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewLogger,
	})
	if err != nil {
		panic(err.Error())
	}
	// DB.AutoMigrate(&model.User{})
}
