package common

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger = logrus.New()

func InitLogger(config *Config) {

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Common.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	// 判断是否输出到文件
	if config.Common.LogToFile {
		// 确保日志目录存在
		err := os.MkdirAll(config.Common.LogDir, 0755)
		if err != nil {
			Logger.Fatalf("Failed to create log directory: %v", err)
		}
		// 配置 lumberjack 日志管理
		Logger.SetOutput(&lumberjack.Logger{
			Filename:   config.Common.LogDir + "/Hedwig.log", // 基础日志文件名
			MaxSize:    config.Common.LogMaxSize,                 // 每个日志文件最大大小（单位 MB）
			MaxAge:     config.Common.LogMaxAge,                  // 日志文件保留天数
			MaxBackups: config.Common.LogMaxCount,                // 最多保留的旧日志文件数
			Compress:   true,                                     // 是否压缩旧日志文件
		})
	} else {
		// 默认输出到终端
		Logger.SetOutput(os.Stdout)
	}
}

/*func InitLogger(level string) {
	Logger.Out = os.Stdout
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	SetLogLevel(level)
}

func SetLogLevel(level logrus.Level) {
	switch level {
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}*/
