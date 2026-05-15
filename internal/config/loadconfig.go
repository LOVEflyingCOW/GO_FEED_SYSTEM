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
	Storage  StorageConfig  `mapstructure:"storage"`
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

type StorageConfig struct {
	UploadDir string `mapstructure:"upload_dir"`
	BaseURL   string `mapstructure:"base_url"`
}

// func Load(filename string) (Config, error) {
// 	v := viper.New()
// 	v.SetConfigFile(filename)
// 	v.SetConfigType("yaml")
// 	v.AutomaticEnv()
// 	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 	if err := v.ReadInConfig(); err != nil {
// 		return Config{}, err
// 	}

// 	var cfg Config
// 	return cfg, v.Unmarshal(&cfg)
// }

func LoadLocalDev(filename string) (Config, bool, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	defaults := DefaultLocalConfig()
	v.SetDefault("server.port", defaults.Server.Port)
	v.SetDefault("database.host", defaults.Database.Host)
	v.SetDefault("database.port", defaults.Database.Port)
	v.SetDefault("database.user", defaults.Database.User)
	v.SetDefault("database.password", defaults.Database.Password)
	v.SetDefault("database.dbname", defaults.Database.DBName)
	v.SetDefault("redis.host", defaults.Redis.Host)
	v.SetDefault("redis.port", defaults.Redis.Port)
	v.SetDefault("redis.password", defaults.Redis.Password)
	v.SetDefault("redis.db", defaults.Redis.DB)
	v.SetDefault("rabbitmq.host", defaults.RabbitMQ.Host)
	v.SetDefault("rabbitmq.port", defaults.RabbitMQ.Port)
	v.SetDefault("rabbitmq.username", defaults.RabbitMQ.Username)
	v.SetDefault("rabbitmq.password", defaults.RabbitMQ.Password)
	v.SetDefault("storage.upload_dir", defaults.Storage.UploadDir)
	v.SetDefault("storage.base_url", defaults.Storage.BaseURL)

	usedDefault := false
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			usedDefault = true
		} else {
			return Config{}, usedDefault, err
		}
	}

	var cfg Config
	return cfg, usedDefault, v.Unmarshal(&cfg)
}

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
		Storage: StorageConfig{
			UploadDir: "./uploads",
			BaseURL:   "http://localhost:8081/uploads",
		},
	}
}
