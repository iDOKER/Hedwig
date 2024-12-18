package common

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mode string `yaml:"mode"`

	Receiver struct {
		ListenAddress string `yaml:"listen_address"`
		OutputDir     string `yaml:"output_dir"`
		HeaderKey     string `yaml:"header_key"`
		HeaderValue   string `yaml:"header_value"`
		EncryptToken  string `yaml:"encrypt_token"`
	} `yaml:"receiver"`

	Sender struct {
		WatchDir      string `yaml:"watch_dir"`
		TargetURL     string `yaml:"target_url"`
		Timeout       int    `yaml:"timeout"`
		RetryCount    int    `yaml:"retry_count"`
		RetryInterval int    `yaml:"retry_interval"`
		SSLVerify     bool   `yaml:"ssl_verify"`
		HeaderKey     string `yaml:"header_key"`
		HeaderValue   string `yaml:"header_value"`
		EncryptToken  string `yaml:"encrypt_token"`
	} `yaml:"sender"`

	Common struct {
		FilePrefix            string `yaml:"file_prefix"`
		FileSuffix            string `yaml:"file_suffix"`
		LogLevel              string `yaml:"log_level"`
		LogToFile             bool   `yaml:"log_to_file"`
		LogDir                string `yaml:"log_dir"`
		LogMaxSize            int    `yaml:"log_max_size"`
		LogMaxAge             int    `yaml:"log_max_age"`
		LogMaxCount           int    `yaml:"log_max_count"`
		DataBackupEnable      bool   `yaml:"data_backup_enable"`
		DataBackupDir         string `yaml:"data_backup_dir"`
		DataBackupMaxAge      int    `yaml:"data_backup_max_age"`
		DataBackupMaxInterval int    `yaml:"data_backup_max_interval"`
	} `yaml:"common"`
}

func LoadConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		Logger.Errorf("Failed to load config file: %s, error: %v\n", path, err)
		log.Fatalf("Failed to read config file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		Logger.Errorf("Failed to parse config file: %s, error: %v\n", path, err)
		log.Fatalf("Failed to parse config file: %v", err)
	}

	return &config
}
