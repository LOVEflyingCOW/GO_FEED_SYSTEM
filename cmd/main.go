package main

import (
	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/db"
	"log"
	"os"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	log.Printf("Loading config from %s", configPath)
	cfg, usedDefault, err := config.LoadLocalDev(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if usedDefault {
		log.Printf("Config File %s not found, using default local config", configPath)
	} else {
		log.Printf("Config loaded from file: %s", configPath)
	}

	sqlDB, err := db.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.CloseDB(sqlDB)
	log.Printf("Database connected")

	log.Printf("Server is running on port %d", cfg.Server.Port)
	log.Fatal("Server started") // 暂时不真正启动，等路由完成后启用
}
