package common

import (
	"os"
	"path/filepath"
	"time"
)

// StartFileCleaner 启动一个定时任务，定期清理指定目录的旧文件
func StartFileCleaner(directory string, retentionDays int, intervalMinutes int) {
	if retentionDays <= 0 {
		Logger.Errorf("Retention days not set or invalid; skipping file cleaning.")
		return
	}

	if intervalMinutes <= 0 {
		Logger.Errorf("Interval minutes not set or invalid; skipping file cleaning.")
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
		defer ticker.Stop()

		Logger.Infof("File cleaner started. Retention days: %d, Interval: %d minutes. Directory: %s \n", retentionDays, intervalMinutes, directory)

		for range ticker.C {
			Logger.Infof("Starting scheduled file cleaning...[%s]", directory)
			CleanOldFiles(directory, retentionDays)
		}
	}()
}

func CleanOldFiles(directory string, retentionDays int) {
	// 如果未设置保留天数，直接返回
	if retentionDays <= 0 {
		return
	}

	// 计算文件的最早保留时间
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	// 遍历目录中的所有文件
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Logger.Errorf("Failed to access file: %s, error: %v\n", path, err)
			return nil
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 如果文件的修改时间早于保留时间，删除文件
		if info.ModTime().Before(cutoffTime) {
			Logger.Infof("Deleting old file: %s (last modified: %s)\n", path, info.ModTime())
			err := os.Remove(path)
			if err != nil {
				Logger.Errorf("Failed to delete file: %s, error: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		Logger.Errorf("Failed to clean old files in directory: %s, error: %v\n", directory, err)
	}
}
