package main

import (
	"Hedwig/internal/common"
	"Hedwig/internal/receiver"
	"Hedwig/internal/sender"
)

func main() {
	// 加载配置
	config := common.LoadConfig("./config/config.yaml")

	// 初始化日志
	common.InitLogger(config)
	common.Logger.Info("Program started")

	// 清理旧的文件
	if config.Common.DataBackupEnable && config.Common.DataBackupMaxAge > 0 {
		common.StartFileCleaner(config.Common.DataBackupDir, config.Common.DataBackupMaxAge, config.Common.DataBackupMaxInterval)
	} else {
		common.Logger.Infof("Data backup is %v and data_backup_max_age <= 0, skipping file cleaning.", config.Common.DataBackupEnable)
	}

	// 根据模式运行接收端或发送端逻辑
	switch config.Mode {
	case "receiver":
		common.Logger.Infof("Running in Receiver mode, log level: %s", config.Common.LogLevel)
		err := receiver.StartServer(config.Receiver.ListenAddress, config.Receiver.OutputDir, config.Common.FilePrefix, config.Common.FileSuffix, config.Receiver.HeaderKey, config.Receiver.HeaderValue, config.Receiver.EncryptToken, config.Common.DataBackupEnable, config.Common.DataBackupDir)
		if err != nil {
			common.Logger.Infof("Receiver failed: %v", err)
		}

	case "sender":
		common.Logger.Infof("Running in Sender mode, log level: %s", config.Common.LogLevel)
		err := sender.WatchAndForward(config.Sender.WatchDir, config.Common.FilePrefix, config.Common.FileSuffix, config.Sender.TargetURL, config.Sender.Timeout, config.Sender.RetryCount, config.Sender.RetryInterval, config.Sender.SSLVerify, config.Sender.HeaderKey, config.Sender.HeaderValue, config.Sender.EncryptToken, config.Common.DataBackupEnable, config.Common.DataBackupDir)
		if err != nil {
			common.Logger.Infof("Sender failed: %v", err)
		}

	default:
		common.Logger.Errorf("Invalid mode: %s, expected 'receiver' or 'sender'", config.Mode)
	}
}
