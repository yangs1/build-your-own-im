package initialize

import (
	"fmt"
	"gim/global"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitDB() {
	// MySQL 配置信息
	username := "root"   // 账号
	password := "123456" // 密码
	host := "127.0.0.1"  // 地址
	port := 3306         // 端口
	DBname := "gim"      // 数据库名称
	timeout := "10s"     // 连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
	// Open 连接
	//写sql语句配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	var err error

	//将获取到的连接赋值到global.DB
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, //打印sql日志
	})
	if err != nil {
		panic(err)
	}
}

func InitRedis() {
	host := "127.0.0.1" // 地址
	port := "6379"
	opt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port), // redis地址
		Password: "",                               // redis密码，没有则留空
		DB:       10,                               // 默认数据库，默认是0
	}
	global.RedisDB = redis.NewClient(&opt)
}
