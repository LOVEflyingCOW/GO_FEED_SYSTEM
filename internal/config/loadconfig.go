package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func Load(filename string) (Config, error) {
	v := viper.New()

	v.SetConfigFile(filename)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	//用 strings.NewReplacer 替换环境变量中的点为下划线
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// LoadLocalDev 加载本地开发配置（如果配置文件不存在则使用默认配置）
func LoadLocalDev(filename string) (Config, bool, error) {
	cfg, err := Load(filename)
	if err == nil {
		return cfg, false, nil
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return DefaultLocalConfig(), true, nil
	}
	return Config{}, false, err
}

// 环境变量 > 配置文件 > 默认值
func DefaultLocalConfig() Config {
	return Config{
		Server: ServerConfig{
			Port: 8081,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "password",
			DBName:   "feedsystem",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
		RabbitMQ: RabbitMQConfig{
			Host:     "localhost",
			Port:     5672,
			Username: "guest",
			Password: "guest",
		},
	}
}
