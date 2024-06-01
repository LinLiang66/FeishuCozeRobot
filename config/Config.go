package config

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	AppName     string      `json:"app_name"`
	AppModel    string      `json:"app_model"`
	AppHost     string      `json:"app_host"`
	AppPort     string      `json:"app_port"`
	RedisConfig RedisConfig `json:"redis_config"`
	BotID       string      `json:"bot_id"`
	BotToken    string      `json:"bot_token"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

// GetConfig 获取配置，外部使用"config.GetConfig()"调用
func GetConfig() *Config {
	return cfg
}

// 存储配置的全局对象
var cfg *Config = nil

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path) //读取文件
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader) //解析json
	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func ChangeConfig(Jsonstr string) (*Config, error) {

	err := json.Unmarshal([]byte(Jsonstr), &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
