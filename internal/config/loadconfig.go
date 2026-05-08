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

// type Config struct {
// 	Server   ServerConfig   `yaml:"server"`
// 	Database DatabaseConfig `yaml:"database"`
// 	Redis    RedisConfig    `yaml:"redis"`
// 	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
// }

// type ServerConfig struct {
// 	Port int `yaml:"port"`
// }

// type DatabaseConfig struct {
// 	Host     string `yaml:"host"`
// 	Port     int    `yaml:"port"`
// 	User     string `yaml:"user"`
// 	Password string `yaml:"password"`
// 	DBName   string `yaml:"dbname"`
// }

// type RedisConfig struct {
// 	Host     string `yaml:"host"`
// 	Port     int    `yaml:"port"`
// 	Password string `yaml:"password"`
// 	DB       int    `yaml:"db"`
// }

// type RabbitMQConfig struct {
// 	Host     string `yaml:"host"`
// 	Port     int    `yaml:"port"`
// 	Username string `yaml:"username"`
// 	Password string `yaml:"password"`
// }

// // Load 加载配置文件
// func Load(filename string) (Config, error) {
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		return Config{}, fmt.Errorf("failed to read config file: %w", err)
// 	}

// 	var cfg Config
// 	if err := yaml.Unmarshal(data, &cfg); err != nil {
// 		return Config{}, fmt.Errorf("parse %s: %w", filename, err)
// 	}

// 	ApplyEnvOverrides(&cfg)
// 	return cfg, nil
// }

// // ApplyEnvOverrides 环境变量覆盖配置（支持 Docker/K8s 部署）
// func ApplyEnvOverrides(cfg *Config) {
// 	if cfg == nil {
// 		return
// 	}
// 	if v := os.Getenv("SERVER_PORT"); v != "" {
// 		if port, err := strconv.Atoi(v); err == nil {
// 			cfg.Server.Port = port
// 		}
// 	}
// 	if v := os.Getenv("MYSQL_HOST"); v != "" {
// 		cfg.Database.Host = v
// 	}
// 	if v := os.Getenv("MYSQL_PORT"); v != "" {
// 		if port, err := strconv.Atoi(v); err == nil {
// 			cfg.Database.Port = port
// 		}
// 	}
// 	if v := os.Getenv("MYSQL_USER"); v != "" {
// 		cfg.Database.User = v
// 	}
// 	if v := os.Getenv("MYSQL_ROOT_PASSWORD"); v != "" {
// 		cfg.Database.Password = v
// 	}
// 	if v := os.Getenv("MYSQL_DATABASE"); v != "" {
// 		cfg.Database.DBName = v
// 	}
// 	if v := os.Getenv("REDIS_HOST"); v != "" {
// 		cfg.Redis.Host = v
// 	}
// 	if v := os.Getenv("REDIS_PORT"); v != "" {
// 		if port, err := strconv.Atoi(v); err == nil {
// 			cfg.Redis.Port = port
// 		}
// 	}
// 	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
// 		cfg.Redis.Password = v
// 	}
// 	if v := os.Getenv("RABBITMQ_HOST"); v != "" {
// 		cfg.RabbitMQ.Host = v
// 	}
// 	if v := os.Getenv("RABBITMQ_PORT"); v != "" {
// 		if port, err := strconv.Atoi(v); err == nil {
// 			cfg.RabbitMQ.Port = port
// 		}
// 	}
// 	if v := os.Getenv("RABBITMQ_USER"); v != "" {
// 		cfg.RabbitMQ.Username = v
// 	}
// 	if v := os.Getenv("RABBITMQ_PASS"); v != "" {
// 		cfg.RabbitMQ.Password = v
// 	}
// }

// // LoadLocalDev 加载本地开发配置（如果配置文件不存在则使用默认配置）
// func LoadLocalDev(filename string) (Config, bool, error) {
// 	cfg, err := Load(filename)
// 	if err == nil {
// 		return cfg, false, nil
// 	}
// 	if errors.Is(err, os.ErrNotExist) {
// 		return DefaultLocalConfig(), true, nil
// 	}
// 	return Config{}, false, err
// }

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
