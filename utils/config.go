package utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Es EsConfig `json:"es" yaml:"ElasticSearch"`
}

type EsConfig struct {
	Address  string `json:"address" yaml:"Address"`
	Username string `json:"username" yaml:"Username"`
	Password string `json:"password" yaml:"Password"`
}

// 获取配置文件
func GetConfig() *Config {
	var config *Config
	configFilePath := filepath.Join("config.yaml")
	// 读取本地配置 yaml 文件
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		panic(err)
	}
	return config
}