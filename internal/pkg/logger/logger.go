package logger

import (
	"fmt"
	"log"
)

// Warn 记录警告级别日志
func Warn(module, message string, args ...interface{}) {
	log.Printf("[WARN] [%s] %s", module, fmt.Sprintf(message, args...))
}

// Error 记录错误级别日志
func Error(module, message string, args ...interface{}) {
	log.Printf("[ERROR] [%s] %s", module, fmt.Sprintf(message, args...))
}

// Info 记录信息级别日志
func Info(module, message string, args ...interface{}) {
	log.Printf("[INFO] [%s] %s", module, fmt.Sprintf(message, args...))
}

// Debug 记录调试级别日志
func Debug(module, message string, args ...interface{}) {
	log.Printf("[DEBUG] [%s] %s", module, fmt.Sprintf(message, args...))
}
